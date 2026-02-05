# Wacdo

Une plateforme de gestion de commandes pour Wacdo, une enseigne fictive de restauration rapide, développée en Go avec le framework Gin.

## Installation et configuration

### Prérequis

- **Go** (>= 1.25.0)
- **PostgresL**

### Configuration des variables d'environnement

Créer un fichier `.env` à la racine du projet, en reprenant le contenu du fichier `.env.dist`, et en le personnalisant avec vos informations.

### Lancement de l'application

```bash
go run main.go
```

Le serveur démarrera par défaut sur `http://localhost:8080`.

## Documentation

### Swagger

Une fois le serveur lancé, vous pouvez accéder à la documentation Swagger à l'adresse suivante :
`http://localhost:8080/swagger/index.html`

### Postman

Vous pouvez importer la collection Postman en utilisant le fichier `postman_collection.json` placé à la racine du projet.