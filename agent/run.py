import nbformat
from nbconvert.preprocessors import ExecutePreprocessor
from nbconvert import HTMLExporter
from nbparameterise import extract_parameters, parameter_values, replace_definitions


def notebook_to_html(notebook_path, parameters):
    """Run all cells and convert to static html
    `notebook_path`: Path to notebook 
    `parameters`: Dict of param keys to values    
    """
    with open(notebook_path) as f:
        nb = nbformat.read(f, as_version=4)
    orig_parameters = extract_parameters(nb)
    replaced_params = parameter_values(orig_parameters, **parameters)
    parametrised_nb = replace_definitions(nb, replaced_params)
    ep = ExecutePreprocessor(timeout=600, kernel_name='python3')
    ep.preprocess(parametrised_nb, {})
    html_exporter = HTMLExporter()
    html_exporter.template_file = 'template.tpl'
    (body, resources) = html_exporter.from_notebook_node(parametrised_nb)
    return body




#body = notebook_to_html('examples/Basic.ipynb', {'client': 'AGG', 'iteration': 2})

#with open('output.html', 'wt') as f:
#    f.write(body)

