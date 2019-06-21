package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

type pathUrl struct{
	Path 	string `yaml:"path"`
	URL 	string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
//It is like middlevare
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if prov, ok := pathsToUrls[path]; ok{
			http.Redirect(w,r,prov,http.StatusNotFound)
			return
		}
		fallback.ServeHTTP(w,r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	//from data receive url
	parseURL,err := parseYAML(yaml)

	//check the error
	if err != nil{
		return nil, err
	}

	//create slice where we will loking for url
	pathToURl := buildMap(parseURL)

	//call MapHandler for correct redirect
	return MapHandler(pathToURl,fallback),nil
}

//create func for building the slice
func buildMap(pathUrls []pathUrl) map[string]string{
	//create slice
	pathToUrls := make(map[string]string)

	//fullfill the data from url
	for _,pu := range pathUrls{
		pathToUrls[pu.Path] = pu.URL
	}

	//return fullfil slice
	return pathToUrls
}

// data come as yaml
func parseYAML(data []byte)([]pathUrl,error){
	//variable from yaml to our struct type
	var pathUrls [] pathUrl

	//unmarshal data to our structure
	err := yaml.Unmarshal(data,&pathUrls)
	if err != nil{
		return nil,err
	}

	//return data and error
	return pathUrls, nil
}