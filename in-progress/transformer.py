import yaml
import json
import getopt
import sys
import io
from typing import Optional

# Originally, this program would take a JSON file of an array of elements, capitalise them and write them to another file
# then we decided to add YAML support
# then somebody suggested it would be great to be able to take STDIN and emit to STDOUT, so we did that too
# Then somebody asked for lower-case as well. No problem!

# pipenv run flake8 --max-cognitive-complexity=3 transformer.py


def extract_format_from_extension(filename):
    if filename.endswith(".json"):
        return "JSON"
    if filename.endswith(".yml") or filename.endswith(".yaml"):
        return "YAML"


def infer_input_format(input_format: Optional[str], filename: str) -> str:
    if input_format:
        return input_format
    extracted_format = extract_format_from_extension(filename)
    if extracted_format:
        return extracted_format
    raise Exception("Unsupported input format, must be yaml or json")


def apply_transformation(stuff, transformation):
    more_stuff = []

    if transformation == "CAPITALISE":
        for element in stuff:
            more_stuff.append(element.upper())
    elif transformation == "DECAPITALISE":
        for element in stuff:
            more_stuff.append(element.lower())

    return more_stuff


def transform(filename, stream, input_format, output_format, transformation, output_file, output_stream):
    stuff = []
    raw_stuff = ''
    if filename is not None:
        with open(filename) as fp:
            input_format = infer_input_format(input_format, filename)
            if input_format == 'YAML':
                stuff = yaml.safe_load(fp)
            else:
                raw_stuff = fp.read()
                stuff = json.loads(raw_stuff)
    else:
        raw_stuff = stream.read()
        if input_format is None:
            raise Exception("No input format was specified")
        if input_format == "YAML":
            stuff = yaml.load(stuff)
        elif input_format == "JSON":
            stuff = json.load(raw_stuff)

    if output_format is None:
        extracted_format = extract_format_from_extension(filename)
        if extracted_format:
            output_format = extracted_format
        else:
            output_format = input_format

    more_stuff = apply_transformation(stuff, transformation)

    if output_format == "JSON":
        if output_file is not None:
            json.dump(more_stuff, open(output_file, "w"))
            # TODO: shit, how to close this open()?
        elif output_stream is not None:
            json.dump(more_stuff, output_stream)
        else:
            raise Exception("No output stream is specified")

    if output_format == "YAML":
        if output_stream is not None:
            yaml.dump(more_stuff, output_stream)
        else:
            with open(output_file, "w") as fp:
                yaml.dump(more_stuff, fp)

    return


def main():
    options, arguments = getopt.getopt(
        sys.argv[1:],                      # Arguments
        'i:o:f:F:ul',                            # Short option definitions
        ["input", "output", "format", "output_format", "upper", "lower"])  # Long option definitions
    filename = None
    input_stream = None
    input_format = "JSON"
    output_file = None
    output_stream = None
    output_format = None
    transformation = "CAPITALISE"
    for o, a in options:
        if o in ("-i", "--input"):
            if a == "-":
                input_stream = sys.stdin
            else:
                filename = a
        if o in ("-o", "--output"):
            if a == "-":
                output_stream = sys.stdout
            else:
                output_file = a
        if o in ("-f", '--format'):
            input_format = a
        if o in ("-F", '--output_format'):
            output_format = a
        if o in ("-u", "--upper"):
            transformation = 'CAPITALISE'
        if o in ("-l", "--lower"):
            transformation = "DECAPITALISE"
    transform(filename, input_stream, input_format,
              output_format, transformation, output_file, output_stream)


if __name__ == "__main__":
    main()
