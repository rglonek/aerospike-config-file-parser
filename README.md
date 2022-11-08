# Aerospike Configuration File Parser

## Types

```go
// stanza will be of type map[string]nil, map[string]string, map[string]stanza
type stanza map[string]interface{}

// types returned by the stanza.Type() function
type ValueType string

const ValueString = ValueType("string")
const ValueNil = ValueType("nil")
const ValueStanza = ValueType("stanza")
const ValueUnknown = ValueType("unknown")
```

## Functions

```go
// open a file and pass it to Parse
func ParseFile(f string) (s stanza, err error)

// parse given reader and produce a stanza variable with the parse results
func Parse(r io.Reader) (s stanza, err error)

// create a file and pass handler to Write
func (s stanza) WriteFile(f string, prefix string, indent string, sortItems bool) (err error)

// write back the stanza, as defined, to the writer, optionally sorting the items in preferred aerospike sorting order
// prefix will be appended to every line
// indent will be used to indent each line accordingly
func (s stanza) Write(w io.Writer, prefix string, indent string, sortItems bool) (err error)
```

## Example

```go
func main() {
    // parse file into 's'
	s, err := ParseFile("/some/file")
	if err != nil {
		log.Fatal(err)
	}

    // print types of variables
	fmt.Println(s.Type("service"))
	fmt.Println(s.Stanza("service").Type("proto-fd-max"))

    // adjust value
	s.Stanza("service")["proto-fd-max"] = "30000"

    // write back contents of 's' to screen
	err = s.Write(os.Stdout, "", "    ", true)
	if err != nil {
		log.Fatal(err)
	}
}
```