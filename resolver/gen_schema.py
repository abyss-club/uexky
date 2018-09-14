#!/usr/bin/env python3
import os
import importlib

go_header = '''package resolver

// Schema for api
const Schema = `'''

before_query = '''
schema {
    query: Query
    mutation: Mutation
}

type Query {'''

before_mutations = '''}

type Mutation {'''

before_types = '''}
'''

go_footer = '`'


# input: filename
# output: query, mutation, type
def extract_define(filename: str) -> (str, str, str):
    mod_name = filename[:-3]
    schema_mod = importlib.import_module(f'schema.{mod_name}')
    query = getattr(schema_mod, 'queries', '')
    mutation = getattr(schema_mod, 'mutations', '')
    type_def = getattr(schema_mod, 'types', '')
    return query, mutation, type_def


def gen_go_files(queries: str, mutations: str, types: str) -> str:
    return '\n'.join([go_header, before_query] + queries + [before_mutations] +
                     mutations + [before_types] + types + [go_footer])


def gen_schema():
    querys = []
    mutations = []
    types = []
    for filename in os.listdir("./schema"):
        if filename.endswith('.py') and filename != '__init__.py':
            query, mutation, type_def = extract_define(filename)
            querys.append(query)
            mutations.append(mutation)
            types.append(type_def)
    schema_go_file = gen_go_files(querys, mutations, types)
    with open('schema.go', 'w') as f:
        f.write(schema_go_file)


if __name__ == '__main__':
    gen_schema()
