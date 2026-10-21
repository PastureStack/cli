#!/usr/bin/env python3
"""Publish only checked numeric CLI0.6.19 bytes from an exact successful main run.

Default --plan has no network activity. Explicit --execute requires normal
protected/signed merge, exact required checks, complete security/repro gate,
source-bound archive checksums, immutable-release policy, draft asset digest
readback, then publication and anonymous byte-for-byte readback. Never overwrite.
"""
import argparse
import hashlib
import json
from pathlib import Path, PurePosixPath
import re
import shutil
import subprocess
import tarfile
import tempfile
from urllib.request import Request, urlopen

REPO, VERSION, TAG = 'PastureStack/cli', '0.6.19', 'v0.6.19'


def need(ok, code):
    if not ok:
        raise RuntimeError(code)


def command(*args, optional=False):
    result = subprocess.run(args, capture_output=True, timeout=180)
    need(optional or result.returncode == 0, 'RELEASE_COMMAND_FAILED')
    return result


def api(route):
    return json.loads(command('gh', 'api', 'repos/' + REPO + '/' + route).stdout)


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def unique(root, name):
    paths = [p for p in root.rglob(name) if p.is_file()]
    need(len(paths) == 1, 'EXACT_CANDIDATE_FILE')
    return paths[0]


def candidate_identity(candidate, source, tree):
    need(unique(candidate, 'source-revision.txt').read_text().strip() == source and
         unique(candidate, 'source-tree.txt').read_text().strip() == tree, 'CANDIDATE_SOURCE_IDENTITY')


def create_draft(source, assets, hashes, notes):
    command('gh', 'release', 'create', TAG, *[str(assets / name) for name in sorted(hashes)], '--repo', REPO,
            '--target', source, '--draft', '--title', 'CLI ' + VERSION, '--notes-file', str(notes))


def find_draft(source):
    # GitHub's by-tag endpoint returns404 until a newly created draft creates
    # its tag at publication. Resolve its unique numeric ID instead.
    matches = [item for item in api('releases?per_page=100') if item['tag_name'] == TAG]
    need(len(matches) == 1, 'EXACT_DRAFT_RELEASE')
    draft = api('releases/' + str(matches[0]['id']))
    need(draft['tag_name'] == TAG and draft['target_commitish'] == source and
         draft['draft'] is True and draft['prerelease'] is False, 'EXACT_DRAFT_SOURCE')
    return draft


def publish_draft(release_id):
    command('gh', 'api', '--method', 'PATCH', 'repos/' + REPO + '/releases/' + str(release_id),
            '-F', 'draft=false', '-F', 'prerelease=false')


def require_published(published, hashes):
    need(published['immutable'] is True and published['draft'] is False and
         published['prerelease'] is False and published['tag_name'] == TAG, 'PUBLISHED_IMMUTABLE_REQUIRED')
    actual = {item['name']: item for item in published['assets']}
    need(len(actual) == len(published['assets']) and set(actual) == set(hashes), 'PUBLISHED_ASSET_SET')
    for name, expected in hashes.items():
        need(actual[name]['digest'] == 'sha256:' + expected, 'PUBLISHED_ASSET_DIGEST')


def checked_source(source, run_id):
    need(api('branches/main')['commit']['sha'] == source, 'MAIN_SOURCE_MOVED')
    protection = api('branches/main/protection')
    need(protection['enforce_admins']['enabled'] and protection['required_status_checks']['strict'] and
         api('branches/main/protection/required_signatures')['enabled'], 'NORMAL_PROTECTION_REQUIRED')
    commit = api('commits/' + source)
    need(commit['commit']['verification']['verified'], 'VERIFIED_SOURCE_SIGNATURE')
    prs = [p for p in api('commits/' + source + '/pulls') if p.get('merged_at') and
           p.get('merge_commit_sha') == source and p['base']['ref'] == 'main']
    need(len(prs) == 1, 'EXACT_NORMAL_MERGED_PR')
    required = protection['required_status_checks']['checks']
    runs = api('commits/' + prs[0]['head']['sha'] + '/check-runs?per_page=100')['check_runs']
    need(required, 'REQUIRED_CHECKS_EMPTY')
    for check in required:
        matches = [r for r in runs if r['name'] == check['context'] and r['app']['id'] == check['app_id']]
        need(len(matches) == 1 and matches[0]['status'] == 'completed' and matches[0]['conclusion'] == 'success',
             'EXACT_REQUIRED_CHECK_NOT_SUCCESS')
    run = api('actions/runs/' + str(run_id))
    need(run['head_sha'] == source and run['head_branch'] == 'main' and run['event'] == 'workflow_dispatch' and
         run['path'] == '.github/workflows/security-release-gate.yml' and
         run['status'] == 'completed' and run['conclusion'] == 'success', 'COMPLETE_EXACT_RELEASE_GATE')
    need(api('immutable-releases')['enabled'] is True, 'IMMUTABLE_RELEASES_REQUIRED')
    return commit['commit']['tree']['sha'], prs[0]['html_url'], run['html_url']


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--plan', action='store_true')
    parser.add_argument('--execute', action='store_true')
    parser.add_argument('--expected-source')
    parser.add_argument('--run-id', type=int)
    parser.add_argument('--receipt', type=Path)
    args = parser.parse_args()
    need(not (args.plan and args.execute), 'EXCLUSIVE_PLAN_EXECUTE')
    if not args.execute:
        print(json.dumps({'status': 'PLAN_ONLY', 'repository': REPO, 'tag': TAG,
                          'source': 'explicit exact signed main source', 'overwrite': False,
                          'requiredChecksBypassed': False, 'anonymousReadback': True}))
        return 0
    source = args.expected_source
    need(source and re.fullmatch('[0-9a-f]{40}', source) and args.run_id and args.receipt and
         not args.receipt.exists(), 'EXPLICIT_RELEASE_INPUTS')
    root = Path(__file__).resolve().parents[1]
    need(command('git', '-C', str(root), 'rev-parse', 'HEAD').stdout.decode().strip() == source and
         not command('git', '-C', str(root), 'status', '--porcelain').stdout, 'CLEAN_EXACT_SOURCE_CHECKOUT')
    tree, merged_pr, workflow = checked_source(source, args.run_id)
    need(not api('git/matching-refs/tags/' + TAG), 'IMMUTABLE_TAG_ALREADY_EXISTS')
    with tempfile.TemporaryDirectory(prefix='pasturestack-cli019-release-') as temporary:
        staging = Path(temporary)
        candidate, evidence, assets = [staging / name for name in ('candidate', 'evidence', 'assets')]
        for directory in (candidate, evidence, assets): directory.mkdir()
        command('gh', 'run', 'download', str(args.run_id), '--repo', REPO, '--name',
                'cli-release-candidate-' + VERSION + '-' + source, '--dir', str(candidate))
        command('gh', 'run', 'download', str(args.run_id), '--repo', REPO, '--name',
                'cli-security-' + source, '--dir', str(evidence))
        candidate_identity(candidate, source, tree)
        summary = json.loads(unique(evidence, 'security-summary.json').read_bytes())
        need(summary['status'] == 'PASS_REVIEWED_VENDOR_POLICY' and not summary['errors'] and
             not summary['blocking'] and summary['unregistered'] == 0, 'CHECKED_SECURITY_SUMMARY')
        manifest = unique(candidate, 'artifacts.sha256').read_text().splitlines()
        for line in manifest:
            expected, relative = line.split(None, 1)
            name = PurePosixPath(relative.strip())
            need(not name.is_absolute() and '..' not in name.parts and re.fullmatch('[0-9a-f]{64}', expected), 'MANIFEST_PATH')
            path = candidate.joinpath(*name.parts)
            need(path.is_file() and path.resolve().is_relative_to(candidate.resolve()) and digest(path) == expected,
                 'PRODUCER_ARTIFACT_HASH')
        archive_paths = [p for p in candidate.rglob('*') if p.is_file() and p.name.startswith('pasturestack-cli-') and
                         p.name.endswith(('.tar.gz', '.tar.xz', '.zip'))]
        need(archive_paths and len({p.name for p in archive_paths}) == len(archive_paths), 'EXACT_NUMERIC_ARCHIVES')
        for path in archive_paths: shutil.copyfile(path, assets / path.name)
        # Server consumes the producer's platform-before-version coordinates.
        for platform, extension in (('linux-amd64', '.tar.gz'), ('darwin-amd64', '.tar.gz'), ('windows-amd64', '.zip')):
            original = assets / ('pasturestack-cli-' + platform + '-' + VERSION + extension)
            need(original.is_file(), 'SERVER_PLATFORM_ASSET_REQUIRED')
        for name in ('security-summary.json', 'source.cdx.json', 'product.cdx.json', 'runtime.cdx.json',
                     'source-security.json', 'product-security.json', 'dapper-image-raw.json', 'runtime-image-security.json',
                     'dapper-dpkg.tsv', 'runtime-dpkg.tsv', 'product-go-version.txt', 'docker-cli-go-version.txt', 'source-locks.sha256',
                     'vendor-pending.json', 'dapper.openvex.json', 'cli018-review-evidence.json'):
            shutil.copyfile(unique(evidence, name), assets / ('cli-' + name))
        binary_archive = assets / ('pasturestack-cli-linux-amd64-' + VERSION + '.tar.gz')
        with tarfile.open(binary_archive, 'r:gz') as tar:
            entries = [m for m in tar.getmembers() if m.isfile() and m.name.endswith('/pasturestack')]
            need(len(entries) == 1 and not entries[0].issym(), 'EXACT_LINUX_MEMBER')
            binary = tar.extractfile(entries[0]).read()
            need(binary.startswith(b'\x7fELF'), 'LINUX_ELF_REQUIRED')
            binary_proof = {'member': entries[0].name, 'bytes': len(binary), 'sha256': hashlib.sha256(binary).hexdigest()}
        identity = {'component': 'cli', 'tag': TAG, 'version': VERSION, 'source': source, 'tree': tree,
                    'normalMerge': True, 'mergedPr': merged_pr, 'workflow': workflow,
                    'requiredChecksBypassed': False, 'linuxBinary': binary_proof,
                    'rawFindingCount': summary['rawFindingCount'], 'trackedRows': summary['trackedRows'],
                    'unregistered': 0, 'trackerIsNotNotAffected': True, 'wholeBE18Pass': False}
        (assets / 'cli-source-producer.json').write_text(json.dumps(identity, indent=2) + '\n', encoding='utf8')
        hashes = {p.name: digest(p) for p in sorted(assets.iterdir())}
        (assets / 'SHA256SUMS').write_text(''.join(value + '  ' + name + '\n' for name, value in hashes.items()), encoding='utf8')
        hashes['SHA256SUMS'] = digest(assets / 'SHA256SUMS')
        checked_source(source, args.run_id)
        create_draft(source, assets, hashes, root / ('RELEASE-' + VERSION + '.md'))
        draft = find_draft(source)
        actual = {item['name']: item for item in draft['assets']}
        need(draft['draft'] and set(actual) == set(hashes), 'DRAFT_ASSET_SET')
        for name, expected in hashes.items():
            need(actual[name]['digest'] == 'sha256:' + expected and actual[name]['size'] == (assets / name).stat().st_size,
                 'DRAFT_ASSET_DIGEST')
        checked_source(source, args.run_id)
        publish_draft(draft['id'])
        published = api('releases/' + str(draft['id']))
        require_published(published, hashes)
        need(api('commits/' + TAG)['sha'] == source, 'PUBLISHED_TAG_SOURCE')
        receipt_assets = []
        for item in published['assets']:
            with urlopen(Request(item['browser_download_url'], headers={'User-Agent': 'PastureStack-release-readback'}), timeout=60) as response:
                data = response.read(item['size'] + 1)
                need(response.status == 200 and len(data) == item['size'] and
                     hashlib.sha256(data).hexdigest() == hashes[item['name']], 'ANONYMOUS_ASSET_BYTES')
            receipt_assets.append({'name': item['name'], 'url': item['browser_download_url'], 'bytes': len(data),
                                   'sha256': hashes[item['name']], 'anonymousReadback': True})
        receipt = {**identity, 'immutable': True, 'draft': False, 'prerelease': False,
                   'assets': receipt_assets, 'publicationUrl': published['html_url']}
        with args.receipt.open('x', encoding='utf8') as handle: json.dump(receipt, handle, indent=2); handle.write('\n')
        print(json.dumps({'status': 'PASS_IMMUTABLE_CLI_RELEASE_ANONYMOUS_READBACK', 'url': published['html_url'],
                          'source': source, 'assets': len(receipt_assets), 'receipt': args.receipt.name}))
    return 0


if __name__ == '__main__':
    try:
        raise SystemExit(main())
    except Exception as error:
        code = str(error) if isinstance(error, RuntimeError) and re.fullmatch('[A-Z_]+', str(error)) else type(error).__name__
        print(json.dumps({'status': 'STOPPED_NO_OVERWRITE_NO_BYPASS', 'code': code}))
        raise SystemExit(1)
