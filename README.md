# Go Wordle

Projet Go en cours de réalisation.

## Suivi de la grille de notation

### CLI en mode TUI

- [X] Ecran d'accueil avec les options possibles
- [X] Ecran d'authentification pour se connecter à son compte
- [X] Ecran de paramètres pour configurer l'application
- [X] Ecrans qui proposent la fonctionnalité principale de l'application
- [ ] Stockage des paramètres dans le dossier personnel de l'utilisateur au format JSON
- [ ] Stockage de l'état de l'application dans le dossier personnel de l'utilisateur au format JSON
- [ ] Import et export de la configuration et de l'état de l'application depuis et vers le serveur
- [X] Interface utilisateur agréable, colorée, dynamique et facile à utiliser

### Serveur d'API REST

- [X] Endpoints d'authentification
- [ ] Endpoints pour récupérer et mettre à jour les données de l'application
- [X] Stockage de données dans une base de données PostgreSQL
- [ ] Endpoints pour importer et exporter la configuration et l'état de l'application depuis et vers le client CLI

### Livraison et déploiement

- [X] Nécessite de cross-compiler le serveur pour tourner sous Linux
- [ ] Héberger le serveur gratuitement
- [ ] Cross-compiler la CLI pour Windows, Linux et MacOS
- [ ] Publier la CLI sur GitHub dans la section Releases du dépôt
- [ ] Répondre aux questions techniques lors de la soutenance

## Release avec GoReleaser

- La configuration GoReleaser est dans [`.goreleaser.yaml`](.goreleaser.yaml).
- Le workflow GitHub Actions se trouve dans [`.github/workflows/release.yml`](.github/workflows/release.yml) et s'exécute sur les tags `v*`.
- Les binaires générés sont `go-wordle-server` pour le backend et `go-wordle` pour le TUI.

## Notes

- État actuel du dépôt: un simple point d'entrée Go qui affiche "Hello, World!".
- Cochez chaque case au fur et à mesure de l'avancement réel du projet.