# Refactoring kata - the transformer

This is the beginnings of a kata used for workshops covering "Separation of responsibility" and refactoring.

## Installation

```sh
make setup
```

This will download all of the Go dependencies, including the `gocognit` binary (for measuring complexity)

## Testing

To run the tests:

```sh
make test
```

To show test coverage:

```sh
make coverage
```

To show complexity:

```sh
make complexity
```

## The exercise

The `transformer` application is a simple application to transform an array of values:

- The input is a file or STDIN with one of these formats: JSON, YAML.
- The currently available transformations are CAPITALISE and DECAPITALISE, which convert the array elements to upper or lower case.
- The output is a file or STDOUT with one of these formats: JSON, YAML.

So, it is possible to transform a JSON array to an upper-case YAML array, and to pipe an input stream into a file, say.

Originally, this application was just intended to capitalise a JSON file. Then we added decapitalisation. Then somebody asked for YAML as well. Then we added a new requirement for STDIN/STDOUT. Being able to convert between JSON and YAML, and vice-versa, was an unexpected bonus.

We have plans to add a third file format - CSV, in the near future, and to support camel-casing ("foo" -> "Foo").

The way the application has evolved has left it rather hard to test. We are getting good (though not perfect) test coverage,
but we are not testing every permutation, and there is concern that some untested permutations might result in unexpected behaviour. The problem is, we have identified at least 19 unwritten tests to cover these permutations. Writing this many tests
is painful, but would need to be done as things stand. If after that we then introduce, say, a new file format, then we expect that this will at least double the number of permutations to be tested. And so forth. This is unsupportable.

Your mission, then, is to work out how to refactor the `Transform()` function so the number of required tests does not grow exponentially.

### How to proceed

The big problem with the code as it stands is that one function is trying to do everything. This results in high cyclomatic complexity - the number of possible paths through the function. To be absolutely confident in the code, there really should be a test for every possible path. Some of the paths are trivial - for example errors in file I/O operations, so long as they are caught in the code, are probably not worth testing.

We suggest you start by identifying individual "concerns" in the function. A concern is one area of responsibility. For example, the capitalisation step is a single concern, separate from any I/O or JSON/YAML marshalling steps. If you can separate these out into their own functions, then you will be able to test all of the edge cases of that concern in isolation from the rest. Then with some clever coding to bring things back together, the main `Transform()` function should itself not need so many tests.

Perhaps as a separate exercise, you could try to write the function from scratch, using test-driven development techniques. If you identify the concerns as you code and separate them straight away, you should see a marked difference in how the tests look. Hopefully, they will look a lot cleaner. The original tests, because they are testing permutations, require a substantial amount of data preparation, file creation, etc.

### Complexity

Another marker for good code is how complex it is to read and understand. Try running `make complexity` on the original. This gives a "cognitive complexity" score. A function that is trivial to read will have a score of 1 or 0. A score of above 10 could be considered to be getting refactor-worthy. Over 50, it looks like a real mess.

Remember, code is written to be read by lesser humans than yourselves!

See https://github.com/uudashr/gocognit for a detailed description of cognitive complexity.
