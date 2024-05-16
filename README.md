#  ATLA Api

## Members
```
Gaaze Verena, 22B030332

```

## API
```
POST /characters: Create new character
GET /characters/:id: Get info about character
PUT /characters/:id: Update info about character 
DELETE /characters/:id: Delete character
```

## DB Structure

```

Table characters {
  id bigserial [primary key]
  name text
  age text
  gender text
  status text
  nation text
  created_at timestamp
  updated_at timestamp
}

Table episodes {
  id bigserial [primary key]
  title text
  air_date date
  created_at timestamp
  updated_at timestamp
}

// many-to-many
Table characters_and_episodes {
  id bigserial [primary key]
  character_id bigserial
  episode_id bigserial
  created_at timestamp
  updated_at timestamp
}

Ref: characters_and_episodes.character_id < characters.id
Ref: characters_and_episodes.episode_id < episodes.id

```