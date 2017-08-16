package core

import (
	"strings"
	"testing"
)

var jsonData = `
{  
   "glossary":{  
      "title":"example glossary",
      "GlossDiv":{  
         "title":"S",
         "GlossList":{
			 "Members": [
				 { "Name": "William"},
				  {"Name": "Michael"}
			 ],
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
   },
   "appendix": [
	   {
		   "title": "My Title",
			"page": 21
	   },
	   {
		   "title": "My Other Title",
			"page": 31
	   },
	   {
		   "title": "My Other Other Title",
			"page": 41
	   }
   ],
   "test": [
	   [
		   {
			   "name": "Test"
		   }
	   ]
   ]
}
`

func TestJsonKey(t *testing.T) {
	m := []byte(jsonData)
	tt := []struct {
		name  string
		key   string
		value string
	}{
		{
			"nested key one level down",
			"glossary.title",
			"example glossary",
		},
		{
			"deep nested key",
			"glossary.GlossDiv.GlossList.GlossEntry.GlossDef.para",
			"A meta-markup language, used to create markup languages such as DocBook.",
		},
		{
			"nested key with array index",
			"glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso[1]",
			"XML",
		},
		{
			"nested key with array index and another key",
			"glossary.GlossDiv.GlossList.Members[1].Name",
			"Michael",
		},
		{
			"root key with array index and nested key",
			"appendix.[1].title",
			"My Other Title",
		},
		{
			"array within array object key",
			"test[0].[0].name",
			"Test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if v := jsonKey(tc.key, m); strings.TrimSpace(tc.value) != strings.TrimSpace(string(trimQuotes(v))) {
				t.Fatalf("expected [key: %s] to have [value: %s] but got [value: %s]", tc.key, tc.value, v)
			}
		})
	}
}
