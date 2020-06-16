# gomeili

TODO:

* [x] Handle additional pages for discourse site (`topic_list.more_topics_url`)
* [x] `go fmt` everything
* [ ] Update README
* [ ] Better comments throughout
* [ ] Templates to shortcode format for accordians/marketplace listing
* [ ] General refactor
  * [ ] separate general `parse` utility?
* [ ] Move self-indexing to marketing site (general, outside this repo)
* [ ] Allow API Docs self-indexing?:
  * [ ] Add local spec support to `openapi`.
* [ ] Clean up modules, what is shared public, what is internal.
* [ ] **TESTS**
* [ ] Clean up: better variable names (i.e. `p` isn't a good name for a `Config` object).
* [ ] Remove unicode characters from text.
  * [ ] Apply this  fix to all text fields. `Questions & Answers` section has this too.
* [ ] Find a better, cleaner way to handle the `Description` field.
* [ ] Judge `config.yaml` field names for obviousness. Rename where appropriate.
* [ ] Use Meilisearch go wrapper to provide a single executable to projects (i.e. fully replace `main.py`)
  * [ ] define host (local v. Platform)
  * [ ] set primary key
  * [ ] set index index name
  * [ ] displayed  attributes (central)
  * [ ] searchable attributes (central)
  * [ ] synonyms (central)
  * [ ] ranking rules (central)
  * [ ] master key check
  * [ ] Create the client
  * [ ] delete old indices
  * [ ] create index
  * [ ] update  synonyms
  * [ ] update settings
  * [ ] add documents


## Sources

```yaml
- name: identifiableName
  type: <type>
  rank: <rank>
  url: "remote_url"
  file: "local_file.ext"
  destination: "local_save_location"
  ignore:
    - "ignored_case"
```

* `name`: Something identifiable for your configuration file.
* `type`:
    * `remote`: Remote pre-built Meilisearch index.
    * `local`: Local pre-build Meilisearch index.
    * `openapi`: OpenAPI 3.0. Using this type will parse the given `url` into an index.
    * `discourse`: Discourse site via its API. Using this type will parse the given `url` into an index.
    * `templates`: Platform.sh templates via template-builder project.
* `rank`: Relevant for multi-site search. `1` given highest priority in results - the greater the value, the less relevant.
* `url`:
* `file`:
* `destination`: Save location of the resource's completed index. If none provided, it will not be written at all.
* `ignore`: Ignores a certain case retrieved from the unparsed data. What is ignored will depend on the resource `type`.
    * `openapi`: Ignores paths that contain tags list. Typically when using ReDoc, this should be  `"NO_INCLUDE"`.
    * `discourse`: Category of posts to be ignored.
    * `templates`: File/template ignored in `template-builder/templates`. Currently just `__init__.py`.
