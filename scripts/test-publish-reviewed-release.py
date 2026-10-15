#!/usr/bin/env python3
"""Offline publisher guard tests; never contacts GitHub or publishes assets."""
import importlib.util
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

SCRIPT = Path(__file__).with_name('publish-reviewed-release.py')
SPEC = importlib.util.spec_from_file_location('publisher', SCRIPT)
publisher = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(publisher)


class PublisherTests(unittest.TestCase):
    def test_create_is_draft_not_prerelease(self):
        with patch.object(publisher, 'command') as command:
            publisher.create_draft('a' * 40, Path('assets'), {'one.zip': 'b' * 64}, Path('notes.md'))
        arguments = command.call_args.args
        self.assertIn('--draft', arguments)
        self.assertNotIn('--prerelease', arguments)
        self.assertIn('assets/one.zip', [str(value).replace('\\', '/') for value in arguments])

    def test_publish_explicitly_clears_both_flags(self):
        with patch.object(publisher, 'command') as command:
            publisher.publish_draft(123)
        self.assertEqual(command.call_args.args,
                         ('gh', 'api', '--method', 'PATCH', 'repos/PastureStack/cli/releases/123',
                          '-F', 'draft=false', '-F', 'prerelease=false'))

    def test_lookup_uses_numeric_id_before_tag_exists(self):
        draft = dict(id=123, tag_name=publisher.TAG, target_commitish='a' * 40,
                     draft=True, prerelease=False)
        with patch.object(publisher, 'api', side_effect=[[draft], draft]) as api:
            self.assertEqual(publisher.find_draft('a' * 40), draft)
        self.assertEqual([call.args[0] for call in api.call_args_list],
                         ['releases?per_page=100', 'releases/123'])

    def test_lookup_rejects_duplicates_and_wrong_source(self):
        draft = dict(id=123, tag_name=publisher.TAG, target_commitish='a' * 40,
                     draft=True, prerelease=False)
        with patch.object(publisher, 'api', return_value=[draft, draft]):
            with self.assertRaisesRegex(RuntimeError, 'EXACT_DRAFT_RELEASE'):
                publisher.find_draft('a' * 40)
        with patch.object(publisher, 'api', side_effect=[[draft], draft]):
            with self.assertRaisesRegex(RuntimeError, 'EXACT_DRAFT_SOURCE'):
                publisher.find_draft('b' * 40)

    def test_published_prerelease_is_rejected(self):
        release = dict(immutable=True, draft=False, prerelease=False, tag_name=publisher.TAG,
                       assets=[dict(name='one.zip', digest='sha256:' + 'b' * 64)])
        publisher.require_published(release, {'one.zip': 'b' * 64})
        for key in ('prerelease', 'draft'):
            with self.assertRaisesRegex(RuntimeError, 'PUBLISHED_IMMUTABLE_REQUIRED'):
                publisher.require_published({**release, key: True}, {'one.zip': 'b' * 64})
        with self.assertRaisesRegex(RuntimeError, 'PUBLISHED_ASSET_DIGEST'):
            publisher.require_published(release, {'one.zip': 'c' * 64})

    def test_missing_or_wrong_tree_is_rejected(self):
        with tempfile.TemporaryDirectory() as temporary:
            candidate = Path(temporary)
            (candidate / 'source-revision.txt').write_text('a' * 40)
            with self.assertRaisesRegex(RuntimeError, 'EXACT_CANDIDATE_FILE'):
                publisher.candidate_identity(candidate, 'a' * 40, 'b' * 40)
            (candidate / 'source-tree.txt').write_text('c' * 40)
            with self.assertRaisesRegex(RuntimeError, 'CANDIDATE_SOURCE_IDENTITY'):
                publisher.candidate_identity(candidate, 'a' * 40, 'b' * 40)
            (candidate / 'source-tree.txt').write_text('b' * 40)
            publisher.candidate_identity(candidate, 'a' * 40, 'b' * 40)

    def test_candidate_upload_preserves_both_identity_files(self):
        workflow = SCRIPT.parents[1] / '.github/workflows/security-release-gate.yml'
        upload = workflow.read_text().split('name: Preserve numeric CLI candidate archives', 1)[1]
        upload = upload.split('name: Clean runner resources', 1)[0]
        for name in ('evidence/source-revision.txt', 'evidence/source-tree.txt', 'evidence/artifacts.sha256'):
            self.assertIn(name, upload)


if __name__ == '__main__':
    unittest.main()
