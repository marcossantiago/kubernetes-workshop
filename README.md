
# Kubernetes Workshop

---

This repository contains the materials for a Kubernetes Workshop and it was forked from this repo [goto-k8s](https://github.com/container-solutions/goto-k8s.git).

Clone this
```
git clone https://github.com/marcossantiago/kubernetes-workshop
```

---

## Outline of the Workshop

* [Setup](./00_setup.md)
* [Introduction](./01_intro.md)
* [Basics](./02_basics.md)
* [A Real Application](./04_real_app.md)
* [Running in Production](./05_productionize.md)
* [Advanced features](./06_advanced.md)
* [Kubernetes Clustering](./06_cluster.md)
* [Cloud Native Monitoring](./07_monitoring.md)
* [Next Step](./08_next_steps.md)

## Run Kubernetes Presentation

Some reveal.js features, like external Markdown and speaker notes, require that presentations run from a local web server. The following instructions will set up such a server as well as all of the development tasks needed to make edits to the reveal.js source code.

1. Install [Node.js](http://nodejs.org/) (4.0.0 or later)

1. Clone the workshop repository

```sh
$ git clone https://github.com/marcossantiago/kubernetes-workshop.git
```

1. Navigate to the kubernetes-workshop folder

```sh
$ cd kubernetes-workshop
```

1. Install dependencies
```sh
$ npm install
```

1. Serve the presentation and monitor source files for changes
```sh
$ npm start
```

1. Open <http://localhost:8000> to view your presentation

   You can change the port by using `npm start -- --port=8001`.