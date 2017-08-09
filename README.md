# JT - JSON Template
Ever wanted to take some json input and just transform is slightly or get out a specific value that is deeply nested. Enter JT!

Say for instance you had the following JSON:
```
{  
   "glossary":{  
      "title":"example glossary",
      "GlossDiv":{  
         "title":"S",
         "GlossList":{  
            "GlossEntry":{  
               "ID":"SGML",
               "SortAs":"SGML",
               "GlossTerm":"Standard GeneralizedMarkup Language",
               "Acronym":"SGML",
               "Abbrev":"ISO 8879:1986",
               "GlossDef":{  
                  "para":"A meta-markup language, used to create markup languages such as DocBook.",
                  "GlossSeeAlso":[  
                     "GML",
                     "XML"
                  ]
               },
               "GlossSee":"markup"
            }
         }
      }
   }
}
```

and you wanted to extract the `GlossTerm` value and format it into a markdown heading. In JT you can do the following:
```
echo '{"glossary":{"title":"example glossary","GlossDiv":{"title":"S","GlossList":{"GlossEntry":{"ID":"SGML","SortAs":"SGML","GlossTerm":"Standard GeneralizedMarkup Language","Acronym":"SGML","Abbrev":"ISO 8879:1986","GlossDef":{"para":"A meta-markup language, used to create markup languages such as DocBook.","GlossSeeAlso":["GML","XML"]},"GlossSee":"markup"}}}}}' | jt -inline-template '### {{ json "GlossDiv.GlossList.GlossEntry.GlossTerm" .glossary | str }}'
```
Which produces the following output:
```
### Standard GeneralizedMarkup Language
```
The full golang templating functionality is available for you to use with `jt`!

In addition to the standard golang template funcs, `jt` adds the following for use:
* `json` - takes the given `key` and uses it to traverse the json structure to produce a value
* `str` - converts whatever is given to it, into a string

## What is happening
Using the template given to `jt` the JSON is parsed and then given to the template. 

If we look at the template line:
```
json "GlossDiv.GlossList.GlossEntry.GlossTerm" .glossary | str
```

The `json` function takes a key string and a json fragment. The key represents the subsequent traversal of the `json` structure to the value you want. Each level of the json the `json` func must traverse is separated by a `.`.

All keys in the `root` of the parsed json can be accessed directly, as can be seen in the exmaple with `.glossary`

## CLI flags
`jt` supports the following cli flags:
* `-input` - read from the given file (default is standard in)
* `-template` - use the given file as the template
* `-inline-template` - use the given string as the template
* `-version` - print version information
