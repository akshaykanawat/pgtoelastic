curl -X PUT "http://localhost:9200/users" -H 'Content-Type: application/json' -d '
{
  "mappings": {
    "properties": {
      "id": { "type": "integer" },
      "name": { "type": "text" },
      "created_at": { "type": "date" },
      "project_ids": { "type": "integer" }
    }
  }
}
'

curl -X PUT "http://localhost:9200/hashtags" -H 'Content-Type: application/json' -d '
{
  "mappings": {
    "properties": {
      "id": { "type": "integer" },
      "name": { "type": "keyword" },
      "project_ids" : { "type": "integer" }
    }
  }
}
'

curl -X PUT "http://localhost:9200/projects" -H 'Content-Type: application/json' -d '
{
  "mappings": {
    "properties": {
      "id": { "type": "integer" },
      "name": { "type": "text" },
      "slug": {
        "type": "text",
        "fields": {
          "fuzzy": {
            "type": "text",
            "analyzer": "fuzzy_analyzer"
          },
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "description": {
        "type": "text",
        "fields": {
          "fuzzy": {
            "type": "text",
            "analyzer": "fuzzy_analyzer"
          }
        }
      },
      "created_at": { "type": "date" },
      "hashtag_ids": { "type": "integer" },
      "user_ids": { "type": "integer" }
    }
  },
  "settings": {
    "index": {
      "analysis": {
        "analyzer": {
          "fuzzy_analyzer": {
            "type": "custom",
            "tokenizer": "standard",
            "filter": ["lowercase", "asciifolding"]
          }
        }
      }
    }
  }
}
'
