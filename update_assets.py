import requests
import svgutils.transform as sg
import cairosvg
from xml.dom import minidom
from os import path
import tempfile
from jinja2 import Template


response = requests.get("https://cocomaterial.com/api/vectors")
data = response.json()
size = 24
for entry in data:
    if not path.exists(path.join("assets", "coco", entry["name"] + ".png")):
        with tempfile.TemporaryDirectory() as directory:
            svgResponse = requests.get(entry["svg"])
            tmpfile = open(path.join(directory, "tmp.svg"), "w")
            tmpfile.write(svgResponse.text)
            tmpfile.close()

            orig = path.join(directory, "tmp.svg")
            fig = sg.fromfile(orig)
            width = float(fig.width[:-2])
            height = float(fig.height[:-2])
            max_size = max(width, height)
            increase_ratio = size / max_size
            new_width = round(width * increase_ratio)
            new_height = round(height * increase_ratio)
            dest = path.join("assets", "coco", entry["name"] + ".png")
            cairosvg.svg2png(url=orig, write_to=dest, parent_width=new_width, parent_height=new_height)

template = Template('''package main

func (p *Plugin) setCocoEntries() {
    p.cocoCategories = map[string][]string{
        {%- for category, entries in coco_categories.items() %}
        "{{category}}": []string{
            {%- for entry in entries %}
            "{{entry}}",
            {%- endfor %}
        },
        {%- endfor %}
    }
    p.cocoEntries = []string{
        {%- for category, entries in coco_categories.items() %}
        {%- for entry in entries %}
        "{{entry}}",
        {%- endfor %}
        {%- endfor %}
    }
}
''')
categories = {}
for entry in data:
    for category in map(lambda x: x.strip(), entry["tags"].split(",")):
        if category not in categories:
            categories[category] = []
        categories[category].append(entry["name"])

f = open(path.join("server", "coco_entries.go"), "w")
f.write(template.render(coco_categories=categories))
f.close()
