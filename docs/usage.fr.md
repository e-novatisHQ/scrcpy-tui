# scrcpy-tui

Lanceur scrcpy en français, livré en exécutable Go autonome pour **Linux amd64**. Les builds Linux arm64 sont disponibles par `make release` mais ne sont pas qualifiés sur une machine arm64. ADB et scrcpy doivent être installés et accessibles dans `PATH`. Aucun téléchargement au lancement. Wi-Fi : l'appareil doit déjà être connecté à ADB ; le lanceur ne configure pas son réseau.

```bash
./bin/scrcpy-tui
```

Appareils et presets sont visibles ensemble. Seuls les appareils en état `device` sont sélectionnables. Le premier appareil disponible est présélectionné ; les dernières sélections sont mémorisées après lancement. Un appareil disparu ne bloque pas la sélection d'un autre.

| Touche | Action |
|---|---|
| Tab, ←, → | Changer de liste |
| ↑, ↓ (ou k, j) | Sélectionner |
| Entrée | Vérifier la cible puis lancer scrcpy |
| Home / End, Page↑ / Page↓ | Premier / dernier, déplacement par page |
| / | Rechercher un appareil par adresse, port ou modèle |
| Échap | Effacer le filtre (ou annuler une vue) |
| ? | Aide complète, ↑/↓ pour défiler |
| l | Dernières sorties scrcpy, ↑/↓ pour défiler |
| r | Rafraîchir ADB |
| n / e | Créer / modifier un preset |
| d | Supprimer le preset ; `o` confirme, toute autre touche annule |
| c | Commande complète ; ↑/↓ pour défiler, Échap pour revenir |
| q / Ctrl+C | Quitter la TUI |

Éditeur : trois champs (nom, description, arguments), Tab/Shift+Tab pour changer de champ, Entrée enregistre, Échap annule. Arguments avec espaces : `--window-title 'Ma TV'`. Les options inconnues sont conservées et scrcpy les valide au lancement. Les options qui remplacent la cible (`--serial`, `-s`, `--tcpip`, leurs abréviations et clusters connus, etc.) sont refusées pour garantir la cohérence de la prévisualisation. Recherche : le filtre s’applique immédiatement, Entrée termine la saisie, Échap efface le filtre. Aucun résultat signifie aucune cible sélectionnée et empêche le lancement. Les états non utilisables restent visibles mais sont ignorés par la navigation. Les actions `n/e/d` portent toujours sur les presets, quel que soit le focus.

Aucun shell, expansion de variable ni substitution de commande n'est exécuté.

La TUI est suspendue pendant scrcpy : ses sorties restent visibles dans l'écran normal, puis le menu revient à sa fermeture. Pour interrompre la session, fermer sa fenêtre ou utiliser Ctrl+C. Les processus de la session sont isolés dans un groupe ; l’arrêt gracieux est borné à deux secondes, puis seuls les processus de ce groupe peuvent être arrêtés de force. SIGTERM pendant la session arrête la session puis quitte la TUI en restaurant le terminal. Le lanceur ne signale aucun processus préexistant ; les options libres transmises à scrcpy conservent leurs propres effets. Une erreur scrcpy est signalée au retour. `l` affiche les dernières sorties standard et erreur, limitées à 32 Kio en mémoire ; aucun journal n’est écrit sur disque. La TUI orchestre une seule session au premier plan ; elle ne prétend pas vérifier que la vidéo est effectivement affichée.

## Configuration

`$XDG_CONFIG_HOME/scrcpy-tui/presets.json` (par défaut `~/.config/scrcpy-tui/presets.json`). `--config /chemin/fichier.json` permet un profil indépendant. Les presets initiaux sont Léger Wi-Fi, Très léger et Qualité. Tous sont modifiables/supprimables. Le fichier n'est créé qu'à l'enregistrement d'un preset ou lors d'un lancement.

Écriture atomique, fichier mode 0600. Un fichier JSON invalide est conservé et des presets par défaut sont chargés avec avertissement ; avant tout enregistrement (y compris celui du dernier appareil au lancement), une copie exacte est créée à côté du fichier sous `presets.json.recovery-*`, en mode 0600. Le nouveau fichier peut ensuite remplacer l’original. Pour restaurer, fermer la TUI, déplacer le fichier courant puis recopier la sauvegarde vers `presets.json`. Les presets individuellement invalides et les noms en double sont ignorés avec avertissement et déclenchent aussi cette copie avant réécriture. Une configuration illisible n’est pas écrasée. Une erreur d'écriture reste affichée. Pour réinitialiser, déplacer le fichier de configuration puis relancer. Ne pas placer de secrets dans les arguments : ils apparaissent dans la commande, la configuration et les processus.

## CLI

Les options sont acceptées avant ou après la sous-commande :

```bash
./bin/scrcpy-tui devices
./bin/scrcpy-tui presets
./bin/scrcpy-tui --device USB123 --preset 'Très léger' preview
./bin/scrcpy-tui --device USB123 --preset 'Léger Wi-Fi' --args "--window-title 'Ma TV'" --yes launch
./bin/scrcpy-tui --version
./bin/scrcpy-tui --help
```

La configuration est un JSON éditable pour gérer les presets sans TUI. `preview` ne lance rien. `launch` revalide l'état ADB. Codes : 0 succès, 1 erreur technique/dépendance/scrcpy, 3 entrée invalide ou absence de terminal, 4 lancement CLI sans `--yes`, 130 interruption du lancement CLI. La sortie normale de la TUI, y compris Ctrl+C, retourne 0.

## Compilation et validation

Go à la version de `.go-version` pour construire une release (minimum source : Go 1.25.0). Les versions de dépendances sont verrouillées dans go.mod/go.sum.

```bash
make build
make tools
make check
make vuln
python3 scripts/pty_test.py
make release
make notices
make fixtures
```

`bin/scrcpy-tui` est le binaire local. `dist/` contient les deux architectures et SHA256SUMS. Pour installer : copier le binaire compatible vers `~/.local/bin/scrcpy-tui`. Les tests utilisent des profils temporaires et des outils factices. Les tests PTY qualifient le clavier, la persistance, les erreurs CLI et la restauration du terminal ; le fonctionnement nominal a été confirmé par l’utilisateur le 18 septembre 2026. Cette confirmation ne qualifie pas l’exécution sur arm64.
