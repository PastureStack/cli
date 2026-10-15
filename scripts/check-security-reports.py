#!/usr/bin/env python3
"""Exact release gate: fixed/High/unregistered block; reviewed unfixed M/L track.

Raw reports are never rewritten. Builder kernel headers may use a reviewed exact
OpenVEX not_affected statement; known applied Ubuntu curl backports use fixed,
never not_affected. Everything else remains a blocking or vendor-pending row.
"""
import argparse
from datetime import date, datetime, timezone
import json
from pathlib import Path
import re

SCOPES = ('source', 'product', 'builder', 'runtime')
FIELDS = ('scope', 'vulnerabilityId', 'purl', 'package', 'installedVersion', 'severity')


def key(row):
    return tuple(row[name] for name in FIELDS)


def assess(reports, tracker, vex, release, today=None):
    today = today or date.today()
    errors, observed, exempted, pending, blocking = [], {}, [], [], []
    if set(reports) != set(SCOPES):
        errors.append('EXACT_FOUR_REPORT_SCOPES')
    statements = {}
    for statement in vex.get('statements', []):
        products = statement.get('products', [])
        if len(products) != 1:
            errors.append('VEX_EXACT_ONE_PRODUCT'); continue
        pair = (statement.get('vulnerability', {}).get('name'), products[0].get('@id'))
        if pair in statements:
            errors.append('VEX_DUPLICATE')
        statements[pair] = statement
    used = set()
    for scope, report in reports.items():
        if report.get('SchemaVersion') != 2 or not isinstance(report.get('Results'), list):
            errors.append('INVALID_SCAN_REPORT_' + scope); continue
        for result in report['Results']:
            if result.get('Secrets'):
                errors.append('SECRET_FINDING_' + scope)
            for item in result.get('Vulnerabilities', []):
                row = {'scope': scope, 'vulnerabilityId': item.get('VulnerabilityID', ''),
                       'purl': item.get('PkgIdentifier', {}).get('PURL', ''),
                       'package': item.get('PkgName', ''), 'installedVersion': item.get('InstalledVersion', ''),
                       'severity': item.get('Severity', ''), 'fixedVersion': item.get('FixedVersion', '')}
                if not all(row[name] for name in FIELDS) or not re.fullmatch(r'(?:CVE-\d{4}-\d+|GO-\d{4}-\d+|GHSA-[\w-]+)', row['vulnerabilityId']):
                    errors.append('FINDING_IDENTITY_' + scope)
                observed[key(row)] = row
                pair = (row['vulnerabilityId'], row['purl'])
                statement = statements.get(pair)
                if statement:
                    impact = statement.get('impact_statement', '')
                    header = (scope == 'builder' and row['package'] == 'linux-libc-dev' and
                              statement.get('status') == 'not_affected' and
                              statement.get('justification') == 'vulnerable_code_not_present' and
                              'headers' in impact and 'kernel image' in impact)
                    curl = (scope == 'builder' and row['package'] in ('curl', 'libcurl4t64', 'libcurl3t64-gnutls') and
                            row['installedVersion'] == '8.18.0-1ubuntu2.8' and
                            statement.get('status') == 'fixed' and
                            'https://ubuntu.com/security/notices/' in impact and
                            '8.18.0-1ubuntu2.8' in impact)
                    if not (header or curl):
                        errors.append('UNSUPPORTED_VEX_BOUNDARY')
                    else:
                        used.add(pair)
                        exempted.append({**row, 'vexStatus': statement['status']})
                        continue
                if row['severity'] not in ('LOW', 'MEDIUM') or row['fixedVersion']:
                    blocking.append(row)
                else:
                    pending.append(row)
    if set(statements) != used:
        errors.append('VEX_EXACT_SET_MISMATCH')
    if tracker.get('schemaVersion') != 1 or tracker.get('release') != release or tracker.get('policy') != {
            'allowedSeverities': ['LOW', 'MEDIUM'], 'requiresEmptyFixedVersion': True}:
        errors.append('TRACKER_POLICY')
    tracked = {}
    for row in tracker.get('findings', []):
        try:
            identity = key(row)
            reviewed = datetime.fromisoformat(row['verifiedAt'].replace('Z', '+00:00')).date()
            valid = (row['scope'] in SCOPES and row['severity'] in ('LOW', 'MEDIUM') and
                     row.get('fixedVersion', '') == '' and row.get('status') == 'vendor-pending' and
                     row['source'] == 'https://ubuntu.com/security/' + row['vulnerabilityId'] and
                     row.get('upstreamVendorStatus') in ('Vulnerable', 'Needs evaluation', 'Deferred', 'Ignored') and
                     date.fromisoformat(row['reviewAfter']) >= today and
                     reviewed <= today and (today - reviewed).days <= 7)
        except (KeyError, TypeError, ValueError):
            errors.append('TRACKER_RECORD'); continue
        if not valid or identity in tracked:
            errors.append('TRACKER_REVIEW_OR_DUPLICATE')
        tracked[identity] = row
    pending_set = {key(row) for row in pending}
    if set(tracked) != pending_set:
        errors.append('TRACKER_EXACT_SET_MISMATCH')
    result = {'release': release, 'status': 'PASS_REVIEWED_VENDOR_POLICY' if not errors and not blocking else 'FAIL_SECURITY_GATE',
              'errors': sorted(set(errors)), 'blocking': blocking, 'vendorPending': sorted(pending, key=key),
              'vexApplied': exempted, 'rawFindingCount': len(observed), 'unregistered': len(pending_set - set(tracked)),
              'trackedRows': len(tracked), 'rawReportsPreserved': True}
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--report', action='append', required=True, help='scope=JSON')
    parser.add_argument('--tracker', type=Path, required=True)
    parser.add_argument('--vex', type=Path, required=True)
    parser.add_argument('--release', required=True)
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    paths = dict(value.split('=', 1) for value in args.report)
    if len(paths) != len(args.report):
        raise SystemExit('DUPLICATE_REPORT_SCOPE')
    reports = {scope: json.loads(Path(path).read_bytes()) for scope, path in paths.items()}
    result = assess(reports, json.loads(args.tracker.read_bytes()), json.loads(args.vex.read_bytes()), args.release,
                    datetime.now(timezone.utc).date())
    args.output.write_text(json.dumps(result, indent=2) + '\n', encoding='utf8')
    print(json.dumps({k: result[k] for k in ('status', 'rawFindingCount', 'trackedRows', 'unregistered', 'errors')}))
    return 0 if result['status'].startswith('PASS_') else 1


if __name__ == '__main__':
    raise SystemExit(main())
