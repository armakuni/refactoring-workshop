import pytest
import json
import yaml

from transformer import transform


def test_json_file_uppercase(tmpdir):
    fp = tmpdir.mkdir("foo").join("input.json")
    fp.write("""[
      "Fred",
      "BOB",
      "arthur"
    ]""")
    output_fp = tmpdir.mkdir("bar").join("output.json")

    transform(str(fp), None, "JSON", "JSON",
              "CAPITALISE", str(output_fp), None)

    stuff = json.loads(output_fp.read())
    assert ['FRED', 'BOB', 'ARTHUR'] == stuff


def test_json_file_lowercase(tmpdir):
    fp = tmpdir.mkdir("foo").join("input.json")
    fp.write("""[
      "Fred",
      "BOB",
      "arthur"
    ]""")
    output_fp = tmpdir.mkdir("bar").join("output.json")

    transform(str(fp), None, "JSON", "JSON",
              "DECAPITALISE", str(output_fp), None)

    stuff = json.loads(output_fp.read())
    assert ['fred', 'bob', 'arthur'] == stuff


def test_yaml_file_lowercase(tmpdir):
    fp = tmpdir.mkdir("foo").join("input.yaml")
    fp.write("""
- Fred
- BOB
- arthur
""")
    output_fp = tmpdir.mkdir("bar").join("output.yaml")

    transform(str(fp), None, "YAML", "YAML",
              "DECAPITALISE", str(output_fp), None)

    stuff = yaml.safe_load(output_fp.read())
    assert ['fred', 'bob', 'arthur'] == stuff


def test_yaml_file_uppercase(tmpdir):
    fp = tmpdir.mkdir("foo").join("input.yaml")
    fp.write("""
- Fred
- BOB
- arthur
""")
    output_fp = tmpdir.mkdir("bar").join("output.yaml")

    transform(str(fp), None, "YAML", "YAML",
              "CAPITALISE", str(output_fp), None)

    stuff = yaml.safe_load(output_fp.read())
    assert ['FRED', 'BOB', 'ARTHUR'] == stuff
