# labtime

A Golang tool for extracting useful information from a Gitlab instance.

## Contents

- [labtime](#labtime)
  - [Contents](#contents)
  - [State of development](#state-of-development)
  - [Software Status](#software-status)
  - [Description](#description)
  - [Configuration](#configuration)
    - [YAML Configuration file](#yaml-configuration-file)
    - [Jinja2 Configuration](#jinja2-configuration)
  - [Roadmap](#roadmap)

## State of development

This application is still under active development. Things may change or break, I will try to reduce the risk of this happening but I cannot commit to non breaking changes at this time.

## Software Status

| Item        | Value |
| ----------- | ----- |
| Version     | 0.0.1 |
| Status      | Alpha |
| Testing     | None  |
| QA          | None  |
| Issues      | None  |
| MRs         | None  |
| Performance | Bad   |


## Description

`labtime` will enable users to pull specific pieces information using API Tokens from the Gitlab API.

## Configuration

### YAML Configuration file
The default configuration file is always written to `config/default/config.yml`, personalised configuration should live under `config/custom/config.yml`.

```yml
instance:
  url: gitlab.somedomain.com
  username: coder1
  accesstoken: "12345"
  project: collection/project
queryconfig:
  timeentrydetectionregex: /time spent/g
  timeentryextractionregex: (?P<month>[0-9]+)(mo)|(?P<week>[0-9]+)(w)|(?P<day>[0-9]+)(d)|(?P<hour>[0-9]+)(h)|(?P<minute>[0-9]+)(m)|(?P<second>[0-9]+)(s)
```

### Jinja2 Configuration

To customise the output of the application edit the Jinja2 template files under `templates`. For instance `templates/IssueOverview.csv.j2` will output a CSV formatted file that will describe the time spent on all stored issues.

```jinja
{% for issue in ctx.ListProjectIssues() %}{{ issue.IID }},{{ ctx.TotalTimeSpent(issue) }}
{% endfor %}
```

This will output a file of the form:

```csv
44,1.000000
43,3.800000
42,0.000000
41,0.200000
40,0.300000
39,0.400000
38,0.000000
37,0.000000
36,0.000000
35,0.000000
34,12.800000
```

The first column is the issue number, the second is the time in booked to that issue in days.

## Roadmap
There are lots of things left to do here, this is a PoC more than an application that should be used. I would be keen to improve performance and documentation. Testing and additional functionality is a must also.

Finally I would like to package this as a `go get`table application.