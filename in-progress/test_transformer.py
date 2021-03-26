import pytest
import json
import yaml

from transformer import transform, infer_input_format, extract_format_from_extension


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

    transform(str(fp), None, None, "YAML",
              "CAPITALISE", str(output_fp), None)

    stuff = yaml.safe_load(output_fp.read())
    assert ['FRED', 'BOB', 'ARTHUR'] == stuff


def test_infer_input_format_None_with_empty_filename_raises_exception():
    with pytest.raises(Exception) as err:
        infer_input_format(None, "")
    assert str(err.value) == "Unsupported input format, must be yaml or json"


def test_infer_input_format_None_with_json_filename_returns_JSON():
    assert infer_input_format(None, "foo.json") == "JSON"


def test_infer_input_format_None_with_yaml_filename_returns_YAML():
    assert infer_input_format(None, "foo.yaml") == "YAML"
    assert infer_input_format(None, "foo.yml") == "YAML"


def test_infer_input_format_None_with_random_filename_raises_exception():
    with pytest.raises(Exception) as err:
        infer_input_format(None, "foo.random") == "JSON"
    assert str(err.value) == "Unsupported input format, must be yaml or json"


def test_infer_input_format_when_provided_format_is_not_None():
    assert infer_input_format("JSON", "foo.yaml") == "JSON"


# def test_infer_input_format_None_with_empty_filename_returns_None():
#     assert infer_input_format(None, "") is None

def test_extract_format_from_extension():
    assert extract_format_from_extension("") == None
    assert extract_format_from_extension("filename.json") == "JSON"
    assert extract_format_from_extension("filename.yml") == "YAML"
    assert extract_format_from_extension("filename.yaml") == "YAML"
    assert extract_format_from_extension("filename") == None
    assert extract_format_from_extension("filename.random") == None
