# gozeppelin

Small utilitaries to help manage Zeppelin notebooks using CLI

## setup
```
    $ go get github.com/ademilly/gozeppelin/...
```

## zeppelinfmt

Format a text file as an input JSON for new note requests.
Text is assumed to be the content of a notebook paragraph.

### usage

```
    $ zeppelinfmt -h
    prints usage
    $ zeppelinfmt -name some_name -filepath path_to_txt_file
    outputs JSON body to stdout
    $ zeppelinfmt -name some_name -filepath path_to_txt_file > some_file.json
    outputs JSON body to some_file.json
    $ cat path_to_txt_file | zeppelin
    ouputs JSON body to stdout
```

## zeppelincli

Client for Zeppelin (0.7.3) Rest API (https://zeppelin.apache.org/docs/0.7.3/rest-api/rest-notebook.html)
Implements:
- login
- new note
