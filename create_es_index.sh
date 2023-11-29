curl -X PUT "http://localhost:9200/projects" -H 'Content-Type: application/json' -d '
{
  "mappings": {
    "properties": {
      "id": { "type": "integer" },
      "name": { "type": "text" },
      "slug": { "type": "text" },
      "description": { "type": "text" },
      "created_at": { "type": "date" }
    }
  }
}
'

curl -X PUT "http://localhost:9200/users" -H 'Content-Type: application/json' -d '
{
  "mappings": {
    "properties": {
      "id": { "type": "integer" },
      "name": { "type": "text" }
    }
  }
}
'

curl -X PUT "http://localhost:9200/projects" -H 'Content-Type: application/json' -d '
{
  "mappings": {
    "properties": {
      "name": {
        "type": "text",
        "fields": {
          "fuzzy": {
            "type": "text",
            "analyzer": "fuzzy_analyzer"
          }
        }
      },
      "slug": {
        "type": "keyword"
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
      "hashtags": {
        "type": "nested",
        "properties": {
          "name": { "type": "keyword" }
        }
      },
      "users": {
        "type": "nested",
        "properties": {
          "name": { "type": "text" },
          "created_at": { "type": "date" }
        }
      }
    }
  },
  "settings": {
    "analysis": {
      "analyzer": {
        "fuzzy_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "edge_ngram"]
        }
      }
    }
  }
}
'

curl -X PUT "http://localhost:9200/hashtags" -H 'Content-Type: application/json' -d '
{
  "mappings": {
    "properties": {
      "id": { "type": "integer" },
      "name": { "type": "text" }
    }
  }
}
'

curl -X PUT "http://localhost:9200/project_hashtags" -H 'Content-Type: application/json' -d '
{
  "mappings": {
    "_parent": { "type": "projects" },
    "properties": {
      "hashtag_id": { "type": "integer" },
      "project_id": { "type": "integer" }
    }
  }
}
'
curl -X PUT "http://localhost:9200/user_projects" -H 'Content-Type: application/json' -d '
{
  "mappings": {
    "_parent": { "type": "projects" },
    "properties": {
      "user_id": { "type": "integer" },
      "project_id": { "type": "integer" }
    }
  }
}
'
