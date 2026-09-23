# Décisions conditionnelles du backlog

État au 23 septembre 2026. Ce document prépare les décisions ; il ne clôture pas
les tâches dont les critères nécessitent une observation, un hôte ou une revue.

## DIST-01 — RPM : différer l'implémentation

La consultation des issues du dépôt pendant cette livraison n'a pas fourni de
demande RPM démontrée. Cela ne prouve pas l'absence de besoins privés. Aucun runner
Fedora/RHEL qualifié n'est identifié. Construire un RPM maintenant ajouterait une
promesse de distribution sans recette représentative.

Déclencheur : demande datée mentionnant distribution/version/architecture et
responsable d'un hôte de recette. Avant un go, arrêter la matrice cible, la source
d'adb/scrcpy, les dépendances, le contenu du paquet et la politique de signature.
Tester installation, upgrade, retrait, conservation du profil et comparaison des
payloads sur chaque famille retenue. Évaluer le coût d'entretien de la CI et le
rollback vers le paquet précédent. Le tar.gz Linux existant reste la distribution
à qualifier pour un besoin ponctuel ; il ne rend pas Fedora officiellement testé.

## DIST-02 — Homebrew : pas de formule à ce stade

Le projet n'annonce pas de support macOS. Une formule ne doit pas précéder le
backend de session testé nativement, les chemins de configuration, les outils
adb/scrcpy, les architectures et une recette matérielle macOS. Un cross-build seul
ne remplit pas ces critères.

Déclencheur : demande macOS explicite, matrice Intel/Apple Silicon décidée et hôte
natif disponible. Après qualification, choisir tap propre ou contribution upstream,
fixer hashes, provenance, mise à jour et retrait ; tester installation, upgrade et
uninstall sans supprimer les profils. Aucune formule ni publication n'est créée
par cette évaluation.

## REL-01 — PR de release automatique : conserver la préparation manuelle

Les versions publiques 0.3.0, 0.3.1 et 0.4.0 consultées ne fournissent pas encore
l'historique de cycles espacés nécessaire à l'évaluation. Le nombre de tags n'est
pas une preuve de stabilité du processus. La CI de tag crée déjà un brouillon et
les attestations ; cette fonctionnalité est distincte d'une PR automatique qui
modifierait VERSION et le changelog.

Pour chaque prochain cycle manuel, enregistrer commit, changements de version,
review, tag signé, exécution CI, digests/SBOM/attestations, qualification, décision
publier/refuser et éventuel rollback. Réévaluer après **trois nouveaux cycles
manuels consécutifs réussis**, seuil proposé à faire adopter par les mainteneurs.
L'automatisation future ouvre une PR révisable ; elle ne contourne ni CODEOWNERS,
ni signature, ni checks, ni décision de publication. Tester un brouillon rejeté,
un échec d'artefact, une relance idempotente et une annulation avant d'activer le bot.

## API-01 — candidat au contrat 1.0, sans annonce de stabilité

Le contrat proposé couvre les sous-commandes `tui`, `devices`, `presets`, `preview`,
`launch`, les options existantes et le lancement CLI explicite `--yes`. Les valeurs
d'options restent intactes même lorsqu'elles ressemblent à des sous-commandes.
L'aperçu n'exécute pas scrcpy. Les options de sélection d'appareil restent protégées.

Codes de sortie actuels à préserver ou à faire migrer explicitement : 0 succès,
1 erreur d'exécution, 3 arguments/configuration de commande invalides, 4 lancement
sans `--yes`, 130 interruption de session. Les textes humains français ne sont pas
une API structurée ; aucun JSON stable n'est promis. Les diagnostics des outils
externes ne sont pas sous contrôle du lanceur.

Le fichier JSON actuel conserve presets et dernières sélections. Avant de figer
1.0, inventorier les champs, leurs valeurs par défaut et les comportements de
lecture/écriture, notamment champs inconnus, profil illisible, JSON invalide,
presets rejetés et copie de récupération. Proposer un champ de version du schéma
et des migrations seulement avec fixtures anciennes et tests aller/retour ; ne
pas prétendre que les anciennes versions préserveront les futurs champs inconnus.
Une migration doit laisser une copie récupérable et refuser les données qu'elle ne
sait pas interpréter sans les écraser.

La compatibilité de session porte sur le processus et ses descendants possédés,
le délai d'arrêt, la restauration du terminal et la conservation des processus
étrangers. Le rendu pixel/parpixel de la TUI, les PID, messages de scrcpy et détails
internes des packages Go ne constituent pas une API publique stable.

Proposition de sortie 1.0 : au moins trois versions multiplateformes stables et
90 jours d'observation après la première Windows qualifiée, sans régression
bloquante non résolue. Ces seuils sont des propositions à approuver, pas une règle
historique acquise. Chaque changement incompatible doit avoir notes de migration,
dépréciation annoncée, preuve de restauration et revue. Ni la période ni les
versions requises ne sont accomplies aujourd'hui.
