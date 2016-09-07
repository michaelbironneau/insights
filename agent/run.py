"""
Insights Agent
--------------
Author: Michael Bironneau <michael.bironneau@openenergi.com>

Runs on the machine with the iPython notebook, polls queue for notebook requests. 

Doesn't need to listen or do stuff on weird ports so can get through firewalls and NAT.

Doesn't need sudo, although running it as a service will.

Command-line arguments:
-----------------------

- server: URL of Insight server, without trailing slash (eg. http://localhost:8080)
- root: Directory where it should look for notebooks (eg. '/home/somebody')
- insecure: Whether to NOT check server certificate when connecting (false by default, only set to true for testing)
- poll-interval: How often to poll the server for new requests (default: 3 seconds)

"""
import nbformat
from nbconvert.preprocessors import ExecutePreprocessor
from nbconvert import HTMLExporter
from nbparameterise import extract_parameters, parameter_values, replace_definitions
import logging
import argparse
import requests
import os.path
import socket
import time

HOSTNAME = socket.getfqdn()

parser = argparse.ArgumentParser()
parser.add_argument("--server", type=str, help="URL of Insight server", default="http://localhost:8080")
parser.add_argument("--root", type=str, help="Root of iPython notebooks", default="")
parser.add_argument("--insecure", dest='insecure', action='store_true')
parser.add_argument("--poll-interval", dest="interval", type=int, help="Polling interval for work queue (seconds)", default=3)
parser.set_defaults(insecure=False)
ARGS = parser.parse_args()

def notebook_to_html(notebook_path, parameters):
    """Run all cells and convert to static html
    `notebook_path`: Path to notebook 
    `parameters`: Dict of param keys to values    
    """
    with open(os.path.join(ARGS.root, notebook_path)) as f:
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


def get_work_item():
    """Get an item from this host's work queue'"""
    url = ARGS.server + "/hosts/{0}/work-items/head".format(HOSTNAME)
    res = requests.get(URL, verify=~ARGS.insecure)
    res.raise_for_status()
    return res.json()


def upload_work(work_item, content):
    """Upload content against work item"""
    URL = ARGS.server + "/hosts/{0}/work-items/{1}/success".format(HOSTNAME, work_item.id)
    res = requests.post(url=URL, data=content, headers={'Content-Type': 'application/octet-stream'}, verify=~ARGS.insecure)
    res.raise_for_status()

def record_exception(work_item, exception):
    """Record exception against work item"""
    URL = ARGS.server + "/hosts/{0}/work-items/{1}/failure".format(HOSTNAME, work_item.id)
    res = requets.post(url=URL, json={"error": exception}, headers={'Content-Type': 'application/json'}, verify=~ARGS.insecure)
    res.raise_for_status()

def process_work_item(item):
    """Process work item"""
    try:
        content = notebook_to_html(item.notebook, item.parameters)
        upload_work(work_item, content)
        logger.debug('Processed work item {0}'.format(item.id))
    except Exception as e:
        logger.exception("Exception processing work item {0}".format(item.id))
        record_exception(item.id, content)

logger = logging.basicConfig(level=0)

while True:
    time.sleep(parser.interval)
    try:
        item = get_work_item()
    except Exception as e:
        logger
    try:
        process_work_item(item)
    except Exception as e:
        logger.exception("Unhandled Exception")