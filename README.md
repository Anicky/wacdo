# Wacdo

Une plateforme de gestion de commandes pour Wacdo, une enseigne fictive de restauration rapide, développée en Go avec le framework Gin.
Ce projet permet aux employés de Wacdo de prendre des commandes, gérer l'avancement de celles-ci, et indiquer leur livraison.

## Fonctionnalités

### Liste des routes

- **Gestion des utilisateurs**
    - Connexion d'un utilisateur
    - Création d'un utilisateur
    - Modification d'un utilisateur
    - Suppression d'un utilisateur
    - Affichage de tous les utilisateurs
    - Affichage d'un utilisateur
- **Gestion des catégories de produits**
    - Création d'une catégorie de produit
    - Modification d'une catégorie de produit
    - Suppression d'une catégorie de produit
    - Affichage de toutes les catégories de produits
    - Affichage d'une catégorie de produit
- **Gestion des produits**
    - Création d'un produit
    - Modification d'un produit
    - Suppression d'un produit
    - Affichage de tous les produits
    - Affichage d'un produit
- **Gestion des menus**
    - Création d'un menu
    - Modification d'un menu
    - Suppression d'un menu
    - Affichage de tous les menus
    - Affichage d'un menu
- **Gestion des commandes**
    - Création d'une commande
    - Modification d'une commande
    - Modification de l'état d'avancement d'une commande (en cours de préparation, préparée, livrée)
    - Affichage de toutes les commandes
    - Affichage du détail d'une commande

### Rôles utilisateurs

- **Administrateur** (`admin`) : peut effectuer toutes les actions
- **Equipier d'accueil** (`greeter`) : peut prendre les commandes, les modifier, et les livrer
- **Préparateur de commande** (`order_picker`) : peut voir les commandes et les préparer 
- **Manager** (`manager`) : peut voir les commandes, les préparer, et les livrer

## Déploiement de l'application

L'application a été déployée sur Render, à l'adresse suivante : https://wacdo.onrender.com

## Installation et configuration

### Prérequis

- **Go** (>= 1.25.0)
- **PostgresSQL**

### Configuration des variables d'environnement

Créer un fichier `.env` à la racine du projet, en reprenant le contenu du fichier `.env.dist`, et en le personnalisant avec vos informations.

### Lancement de l'application

```bash
go run main.go
```

Le serveur démarrera par défaut sur `http://localhost:8080`.

## Documentation

### Swagger

Une fois le serveur lancé, vous pouvez accéder à la documentation Swagger directement depuis `http://localhost:8080` (redirection automatique vers `http://localhost:8080/swagger/index.html`).

### Postman

Vous pouvez importer la collection Postman en utilisant le fichier `postman_collection.json` placé dans le dossier `docs` du projet.