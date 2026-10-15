#!/usr/bin/env python3
"""Offline policy fixtures, not a vulnerability/runtime scan."""
import copy
from datetime import date
import importlib.util
from pathlib import Path
import unittest

spec = importlib.util.spec_from_file_location('cli_security', Path(__file__).with_name('check-security-reports.py'))
gate = importlib.util.module_from_spec(spec); spec.loader.exec_module(gate)
TODAY = date(2026, 10, 15)


class Policy(unittest.TestCase):
    def setUp(self):
        self.reports = {scope: {'SchemaVersion': 2, 'Results': []} for scope in gate.SCOPES}
        self.tracker = {'schemaVersion': 1, 'release': '0.6.18', 'policy': {
            'allowedSeverities': ['LOW', 'MEDIUM'], 'requiresEmptyFixedVersion': True}, 'findings': []}
        self.vex = {'statements': []}

    def add(self, severity='MEDIUM', fixed='', scope='runtime', package='libpcre2-8-0'):
        row = {'VulnerabilityID': 'CVE-2026-103111', 'Severity': severity, 'PkgName': package,
               'InstalledVersion': '1', 'FixedVersion': fixed,
               'PkgIdentifier': {'PURL': 'pkg:deb/ubuntu/' + package + '@1'}}
        self.reports[scope]['Results'] = [{'Vulnerabilities': [row]}]
        return row

    def track(self, scope='runtime', severity='MEDIUM', package='libpcre2-8-0'):
        row = {'scope': scope, 'vulnerabilityId': 'CVE-2026-103111', 'severity': severity,
               'purl': 'pkg:deb/ubuntu/' + package + '@1', 'package': package, 'installedVersion': '1',
               'fixedVersion': '', 'status': 'vendor-pending', 'upstreamVendorStatus': 'Needs evaluation',
               'source': 'https://ubuntu.com/security/CVE-2026-103111', 'verifiedAt': '2026-10-15T00:00:00Z',
               'reviewAfter': '2026-10-20'}
        self.tracker['findings'].append(row)
        return row

    def result(self):
        return gate.assess(self.reports, self.tracker, self.vex, '0.6.18', TODAY)

    def test_clean_empty_tracker_allowed(self):
        self.assertTrue(self.result()['status'].startswith('PASS_'))

    def test_exact_reviewed_unfixed_medium(self):
        self.add(); self.track()
        self.assertTrue(self.result()['status'].startswith('PASS_'))

    def test_unregistered_blocks(self):
        self.add()
        self.assertEqual(self.result()['unregistered'], 1)

    def test_high_cannot_be_tracked(self):
        self.add('HIGH'); self.track(severity='HIGH')
        self.assertTrue(self.result()['blocking'])
        self.assertIn('TRACKER_REVIEW_OR_DUPLICATE', self.result()['errors'])

    def test_official_fixed_medium_blocks(self):
        self.add(fixed='2'); self.track()
        self.assertTrue(self.result()['blocking'])

    def test_expired_duplicate_and_stale_review_block(self):
        self.add(); row = self.track(); row['reviewAfter'] = '2026-10-14'
        self.tracker['findings'].append(copy.deepcopy(row))
        self.assertIn('TRACKER_REVIEW_OR_DUPLICATE', self.result()['errors'])

    def test_extra_tracker_record_blocks(self):
        self.track()
        self.assertIn('TRACKER_EXACT_SET_MISMATCH', self.result()['errors'])

    def test_secrets_block_without_serializing_match(self):
        self.reports['product']['Results'] = [{'Secrets': [{'Match': 'DUMMY_SECRET_NOT_RETAINED'}]}]
        result = self.result()
        self.assertIn('SECRET_FINDING_product', result['errors'])
        self.assertNotIn('DUMMY_SECRET_NOT_RETAINED', str(result))

    def test_exact_builder_headers_vex_only(self):
        item = self.add('CRITICAL', scope='builder', package='linux-libc-dev')
        self.vex['statements'] = [{'vulnerability': {'name': item['VulnerabilityID']},
            'products': [{'@id': item['PkgIdentifier']['PURL']}], 'status': 'not_affected',
            'justification': 'vulnerable_code_not_present', 'impact_statement': 'Only headers; no kernel image or modules.'}]
        self.assertTrue(self.result()['status'].startswith('PASS_'))
        self.reports['runtime']['Results'] = copy.deepcopy(self.reports['builder']['Results'])
        self.assertTrue(self.result()['blocking'])

    def test_generic_not_affected_is_rejected(self):
        item = self.add('HIGH', package='openssl')
        self.vex['statements'] = [{'vulnerability': {'name': item['VulnerabilityID']},
            'products': [{'@id': item['PkgIdentifier']['PURL']}], 'status': 'not_affected',
            'justification': 'vulnerable_code_not_present', 'impact_statement': 'not shipped'}]
        self.assertIn('UNSUPPORTED_VEX_BOUNDARY', self.result()['errors'])


if __name__ == '__main__':
    unittest.main()
