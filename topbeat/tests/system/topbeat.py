import sys
sys.path.append('../../../libbeat/tests/system')
from beat.beat import TestCase


class BaseTest(TestCase):

    @classmethod
    def setUpClass(self):
        self.beat_name = "topbeat"
        self.build_path = "../../build/system-tests/"
        self.beat_path = "../../topbeat.test"
