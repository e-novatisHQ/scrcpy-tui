# UX-01 — prototype isolé de localisation

Décision proposée : français et anglais ; français par défaut. Ce prototype ne
change pas la CLI, la configuration ni la TUI distribuées. Il ne lance aucun outil
externe. Les langues cibles et le périmètre final restent à valider.

La question de l'écran reste « quel appareil et quel preset puis-je lancer ? ».
Avant implémentation, le wireframe retenu conserve la hiérarchie existante :

```text
scrcpy-tui — prototype FR/EN
Prêt à lancer · Entrée
Appareils                         Presets
> Téléviseur de démonstration      > Léger Wi-Fi
  Non autorisé : téléphone test
Entrée lancer · r rafraîchir · q quitter
```

Traduire les libellés et les états affichés, jamais les identifiants ADB, les noms
de presets utilisateur, les options scrcpy ni les clés JSON. Les raccourcis restent
stables dans les deux langues. La future intégration devra distinguer identifiant
de preset intégré et libellé traduit : aujourd'hui le nom du preset participe à
sa sélection et ne peut être traduit arbitrairement sans migration.

Exécuter `python3 experiments/localization/prototype.py --locale en-US` ou omettre
l'option. Priorité proposée : option explicite, LC_ALL, LC_MESSAGES, LANG, français.
Une locale explicite inconnue est refusée ; un environnement inconnu ou C/POSIX
revient au français. Une variable prioritaire inconnue ne laisse pas une variable
inférieure changer silencieusement la sélection. Les suffixes d'encodage et régions
sont normalisés. `--self-test` vérifie cette politique, la parité des clés et les
placeholders des catalogues.

## Avant intégration dans le produit

Recenser tous les messages CLI/TUI/application, les avertissements de configuration,
l'éditeur, les aides et les erreurs ; les textes des outils tiers restent externes.
Faire relire l'anglais par un locuteur, tester le français sans changement et faire
valider le choix de langue. Ne pas afficher de traduction partielle comme une
interface entièrement localisée.

La recette dans un vrai PTY reste **NOT RUN pour une TUI traduite** : le prototype
ci-dessous est seulement textuel. L'intégration devra couvrir vide, nominal,
chargé, appareil non autorisé, filtre actif, édition, refus/annulation, réussite et
échec, aux largeurs 59/60/89/90/160 et hauteurs 12/18/30. Vérifier accents, caractères
larges, troncature, touches visibles et sélection conservée. Les tests du catalogue
ne remplacent pas cette recette. Aucun écran anglais de production n'est annoncé.
