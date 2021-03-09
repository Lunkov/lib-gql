package gql

import (
  "encoding/json"
  "github.com/graphql-go/graphql"
  "github.com/golang/glog"
)

// https://medium.com/tunaiku-tech/what-is-graphql-and-how-is-it-implemented-in-golang-b2e7649529f1
// https://spec.graphql.org/June2018/
// https://blog.logrocket.com/3-tips-for-implementing-graphql-in-golang/
// https://habr.com/ru/company/ruvds/blog/444346/
// https://blog.eleven-labs.com/en/construct-structure-go-graphql-api/

var fieldsGQL = make(graphql.Fields)

func AppendQuery(index string, f *graphql.Field) {
  fieldsGQL[index] = f
}

func GetSchema() graphql.SchemaConfig {
  return graphql.SchemaConfig{
                                Query: graphql.NewObject(graphql.ObjectConfig{
                                Name: "Query",
                                Fields: fieldsGQL})}
}

func GetSelectedFields(selectionPath []string, resolveParams graphql.ResolveParams) []string {
  fields := resolveParams.Info.FieldASTs
  for _, propName := range selectionPath {
    found := false
    for _, field := range fields {
      if field.Name.Value == propName {
        selections := field.SelectionSet.Selections
        fields = make([]*ast.Field, 0)
        for _, selection := range selections {
          fields = append(fields, selection.(*ast.Field))
        }
        found = true
        break
      }
    }
    if !found {
      return []string{}
    }
  }
  var collect []string
  for _, field := range fields {
    collect = append(collect, field.Name.Value)
  }
  return collect
}

func Query(query_str string) []byte  {
  if glog.V(9) {
    glog.Infof("DBG: Query: %s", query_str)
  }
  
  rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fieldsGQL}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		glog.Errorf("failed to create new schema, error: %v", err)
	}

	paramsGQL := graphql.Params{Schema: schema, RequestString: query_str}
	res := graphql.Do(paramsGQL)
	if len(res.Errors) > 0 {
		glog.Errorf("failed to execute graphql operation, errors: %+v", res.Errors)
	}
	rJSON, _ := json.Marshal(res)  
  return rJSON
}

