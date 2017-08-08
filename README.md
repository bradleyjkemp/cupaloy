# cupaloy [![Build Status](https://travis-ci.org/bradleyjkemp/cupaloy.svg?branch=master)](https://travis-ci.org/bradleyjkemp/cupaloy) [![Coverage Status](https://coveralls.io/repos/github/bradleyjkemp/cupaloy/badge.svg?branch=master)](https://coveralls.io/github/bradleyjkemp/cupaloy?branch=master)
Simple golang snapshot testing: test that your changes don't unexpectedly alter the results of your code.

### Installation
```bash
go get -u github.com/bradleyjkemp/cupaloy
```

### Usage
```golang
func TestExample(t *testing.T) {
  result := someFunction()

  // check that the result is the same as the last time the snapshot was updated
  err := cupaloy.Snapshot(result)
  
  // if the result has changed then an error containing a diff will be returned
  if err != nil {
    t.Fatalf("error: %s", err)
  }
}
```

To update the snapshots simply set the ```UPDATE_SNAPSHOTS``` environment variable and run your tests e.g.
```bash
UPDATE_SNAPSHOTS=true go test ./...
```
