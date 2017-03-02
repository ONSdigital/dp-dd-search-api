dp-dd-search-api
================

The ONS website is currently using Elastic Search version 2.4. As a result the elastic search client is restricted to version 3:
https://github.com/olivere/elastic

### Getting started

##### Elastic search install via brew
* Install Elastic Search `brew install elasticsearch@2.4`
* Ensure the cluster.name property is set to `cluster.name: elasticsearch`.
The configuration file can be found at `/usr/local/etc/elasticsearch/elasticsearch.yml`. For some reason it appended my username onto the end of the default clustername.
* Start Elastic Search service `brew services start elasticsearch@2.4`
* Run it `elasticsearch`

##### Elastic search via dp-compose
The dp-compose project requires the native docker for mac (not docker toolbox)

```
git clone git@github.com:ONSdigital/dp-compose.git
cd dp-compose
./run.sh
```

##### Run the search api
 ```
 make debug
 ```

##### Query the search api
 ```
curl localhost:20051/search\?q\=armed
 ```
 
 or navigate to <http://localhost:20051/search?q=armed>
 
#### Query for suggestion
```
curl localhost:20051/suggest?q=armed
```

Note: The suggest endpoint gives far less information and is optimised for auto-completion and quick look searches.

### Configuration

| Environment variable | Default                                        | Description
| -------------------- | ---------------------------------------------- | ----------------------------------------------------
| BIND_ADDR            | :20051                                         | The port to serve HTTP requests
| ELASTIC_SEARCH_NODES | http://127.0.0.1:9200                          | The Elastic Search node addresses comma separated
| ELASTIC_SEARCH_INDEX | ons                                            | The Elastic Search index to update

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright ©‎ 2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
