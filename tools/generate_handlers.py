#!/usr/bin/python
# -*- coding: UTF-8 -*-


import os

API_HANDLERS_BEGIN = '\t// BEGIN AUTO GENERATED API HANDLERS\n'
API_HANDLERS_END = '\t// END AUTO GENERATED API HANDLERS\n'


def generate_api_handlers():
    generated_lines = []

    files = os.listdir('../handler/api')
    for f in files:
        handler_name = os.path.splitext(f)[0]
        function_name = ''.join([x.capitalize()
                                for x in handler_name.split('_')]) + 'Handler'
        with open('../handler/api/' + f, encoding='utf-8') as f:
            content = f.read()
            if content.find('func ' + function_name + '(') != -1:
                generated_lines.append('\trouter.Handle("/api/%s", api.%s)\n' % (
                    handler_name, function_name))
            else:
                if handler_name != 'errors':
                    print('Warning: %s is not a valid handler, %s not found' %
                          (f.name, function_name))

    with open('../handler/register.go', 'r', encoding='utf-8') as f:
        content = f.read()
        api_handlers_begin = content.find(API_HANDLERS_BEGIN)
        api_handlers_end = content.find(API_HANDLERS_END)
        assert (api_handlers_begin != -1 and api_handlers_end != -1)
    with open('../handler/register.go', 'w', encoding='utf-8') as f:
        f.write(content[:api_handlers_begin+len(API_HANDLERS_BEGIN)])
        f.write(''.join(generated_lines))
        f.write(content[api_handlers_end:])


def main():
    generate_api_handlers()


if __name__ == '__main__':
    os.chdir(os.path.dirname(__file__))
    main()
