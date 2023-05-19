package parser

import (
	"go/ast"
	"go/token"
)

func getInitFunctionForType(typeName, idFieldName string, oneToMany []OneToManyRelationshipField, manyToMany []ManyToManyRelationshipField) *ast.FuncDecl {
	receiverName := "receiver"
	function := &ast.FuncDecl{
		Name: &ast.Ident{Name: InitFunctionName},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{{Name: receiverName}},
					Type:  &ast.StarExpr{X: &ast.Ident{Name: typeName}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Tok: token.ASSIGN,
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   &ast.Ident{Name: receiverName},
							Sel: &ast.Ident{Name: IsInitializedFieldName},
						},
					},
					Rhs: []ast.Expr{
						&ast.Ident{Name: "true"},
					},
				},
			},
		},
		Type: &ast.FuncType{Params: &ast.FieldList{}},
	}

	for _, oneToManyRel := range oneToMany {
		initOneToManyRef := &ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   &ast.Ident{Name: receiverName},
					Sel: &ast.Ident{Name: oneToManyRel.FromFieldName},
				},
			},
			Rhs: []ast.Expr{
				&ast.StarExpr{X: &ast.CallExpr{Fun: &ast.IndexExpr{
					Index: &ast.Ident{Name: oneToManyRel.TypeName},
					X: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "lib"},
						Sel: &ast.Ident{Name: ReferenceNavigationListCtor},
					},
				},
					Args: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "lib"},
								Sel: &ast.Ident{Name: ReferenceNavigationListParam},
							},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key: &ast.Ident{Name: "OwnerId"},
									Value: &ast.SelectorExpr{
										X:   &ast.Ident{Name: receiverName},
										Sel: &ast.Ident{Name: idFieldName},
									},
								},
								&ast.KeyValueExpr{
									Key: &ast.Ident{Name: "OwnerTypeName"},
									Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   &ast.Ident{Name: receiverName},
											Sel: &ast.Ident{Name: NobjectImplementationMethod},
										},
									},
								},
								&ast.KeyValueExpr{
									Key: &ast.Ident{Name: "OtherTypeName"},
									Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											Sel: &ast.Ident{Name: NobjectImplementationMethod},
											X: &ast.ParenExpr{
												X: &ast.StarExpr{
													X: &ast.CallExpr{
														Args: []ast.Expr{&ast.Ident{Name: oneToManyRel.TypeName}},
														Fun:  &ast.Ident{Name: "new"},
													},
												},
											},
										},
									},
								},
								&ast.KeyValueExpr{
									Key: &ast.Ident{Name: "ReferringFieldName"},
									Value: &ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"" + oneToManyRel.FieldName + "\"",
									},
								},
								&ast.KeyValueExpr{
									Key:   &ast.Ident{Name: "IsManyToMany"},
									Value: &ast.Ident{Name: "false"},
								},
							},
						},
					},
				}},
			}}

		function.Body.List = append(function.Body.List, initOneToManyRef)
	}

	for _, manyToManyRel := range manyToMany {
		// PartionKeyName and SortKeyName define the two types
		// used in many-to-many relationship
		// here, the type used in field declaration is different
		// the type different from the owner type
		relationshipType := manyToManyRel.PartionKeyName
		if typeName == manyToManyRel.PartionKeyName {
			relationshipType = manyToManyRel.SortKeyName
		}

		initManyToManyRef := &ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   &ast.Ident{Name: receiverName},
					Sel: &ast.Ident{Name: manyToManyRel.FromFieldName},
				},
			},
			Rhs: []ast.Expr{
				&ast.StarExpr{X: &ast.CallExpr{
					Fun: &ast.IndexExpr{
						Index: &ast.Ident{Name: relationshipType},
						X: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "lib"},
							Sel: &ast.Ident{Name: ReferenceNavigationListCtor},
						},
					},
					Args: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "lib"},
								Sel: &ast.Ident{Name: ReferenceNavigationListParam},
							},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key: &ast.Ident{Name: "OwnerId"},
									Value: &ast.SelectorExpr{
										X:   &ast.Ident{Name: receiverName},
										Sel: &ast.Ident{Name: idFieldName},
									},
								},
								&ast.KeyValueExpr{
									Key: &ast.Ident{Name: "OwnerTypeName"},
									Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   &ast.Ident{Name: receiverName},
											Sel: &ast.Ident{Name: NobjectImplementationMethod},
										},
									},
								},
								&ast.KeyValueExpr{
									Key: &ast.Ident{Name: "OtherTypeName"},
									Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											Sel: &ast.Ident{Name: NobjectImplementationMethod},
											X: &ast.ParenExpr{
												X: &ast.StarExpr{
													X: &ast.CallExpr{
														Args: []ast.Expr{&ast.Ident{Name: relationshipType}},
														Fun:  &ast.Ident{Name: "new"},
													},
												},
											},
										},
									},
								},
								&ast.KeyValueExpr{
									Key: &ast.Ident{Name: "ReferringFieldName"},
									Value: &ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"" + manyToManyRel.FieldName + "\"",
									},
								},
								&ast.KeyValueExpr{
									Key:   &ast.Ident{Name: "IsManyToMany"},
									Value: &ast.Ident{Name: "true"},
								},
							},
						},
					},
				}},
			},

			// []ast.Expr{
			// 	&ast.StarExpr{X: &ast.CallExpr{
			// 		Args: []ast.Expr{
			// 			&ast.SelectorExpr{
			// 				X:   &ast.Ident{Name: receiverName},
			// 				Sel: &ast.Ident{Name: idFieldName},
			// 			},
			// 			&ast.CallExpr{
			// 				Fun: &ast.SelectorExpr{
			// 					X:   &ast.Ident{Name: receiverName},
			// 					Sel: &ast.Ident{Name: NobjectImplementationMethod},
			// 				},
			// 			},
			// 			&ast.BasicLit{
			// 				Kind:  token.STRING,
			// 				Value: "\"\"",
			// 			},
			// 			&ast.Ident{Name: "true"},
			// 		},
			// 		Fun: &ast.IndexExpr{
			// 			Index: &ast.Ident{Name: relationshipType},
			// 			X: &ast.SelectorExpr{
			// 				X:   &ast.Ident{Name: "lib"},
			// 				Sel: &ast.Ident{Name: ReferenceNavigationListCtor},
			// 			},
			// 		},
			// 	}},
			// },
		}

		function.Body.List = append(function.Body.List, initManyToManyRef)
	}

	return function
}
