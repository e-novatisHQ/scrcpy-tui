# Vérification des capacités scrcpy

Avant de lancer un preset intégré dont les arguments sont inchangés, le lanceur
résout `scrcpy` dans PATH et exécute ce même programme avec `--help`. Il vérifie
les déclarations des options nécessaires au preset choisi, ainsi que `-s`.
Renommer le preset ne désactive pas cette vérification. Une option absente produit
un diagnostic indiquant de mettre à jour scrcpy ou d'adapter le preset, avant la
revalidation de l'appareil par ADB. L'aperçu de commande reste hors ligne.

La sonde est relancée à chaque préparation : changer de scrcpy dans PATH ne laisse
pas de cache périmé. Elle dispose de cinq secondes ; la collecte des pipes après
l'arrêt du processus est bornée à deux secondes et la sortie à 1 Mio. Une erreur,
une aide non reconnue ou une sortie excessive empêche le lancement du preset.

Les presets personnalisés, les presets intégrés dont les arguments ont été
modifiés et les arguments libres `--args` conservent leur contrat : scrcpy valide
leurs options et valeurs. Ce contrôle ne valide ni les codecs disponibles sur
l'appareil, ni les combinaisons d'options, ni une plage de versions officiellement
qualifiée. Il reconnaît les déclarations d'options de l'aide scrcpy, y compris les
alias courts et les fins de ligne Windows ; une évolution de ce format peut
nécessiter une adaptation du parseur.

## Preuves et limites

Les tests couvrent tous les presets intégrés, un preset renommé, une option
manquante, les références dans le texte explicatif, CRLF, les échecs et délais de
sonde, la sortie excessive et l'absence de contact ADB après incompatibilité.
Le helper natif permet la même exécution sous Windows et Linux ; les tests PTY
contrôlent aussi le parcours de lancement.

L'aide locale de scrcpy 4.1 a été inspectée sans lancer d'appareil. Son SHA-256
était `3b09a9244614a49f586fdf3df35b8322b9c696f2bf6f60fcbda016f75a68c750`.
Cette inspection confirme le format attendu, pas une qualification matérielle
ni la compatibilité de toutes les versions scrcpy. La matrice QUAL-01 devra
enregistrer séparément version, OS, appareil, firmware, transport et résultat.
