## LFL Playoff Resolver

## Disclamer
Le code a ete fait rapidement, il y a sans doutes des bugs. Si vous en voyez, ou que vous voulez aider a ameliorer le package n'hesitez pas a me contacter

## Usage
Les resultat de la league sont dans le fichier `cmd/resolve/planning-reel.json`

Calculer les resultats: `cmd/resolve go run main.go`
Tous les scenarios possibles sont calcules.

## Tiebreak 
Le script supporte les tiebreak suivants:
- Le Head2Head
- Le nombre de victoires en match retours

Si le nombre de victoires en match retours ne permet pas de departager les equipes, la resolution est aleatoire.
Le nombre de fois ou le tiebreak ne peut pas etre resolut est donne dans le resultat.

## Results
On peut voir des petits changement avec les tiebreaks non geres, mais les stats finales sont toujours assez similaires.

### Exemple a la fin du jour 16 de la LFL:
Scenarios compute: 99% (1023)Number of possible scenarios: 1024
- Team GO: 99%
- Team LDLC: 98%
- Team AEG: 91%
- Team BDS: 81%
- Team VIT: 80%
- Team GW: 76%
- Team SLY: 54%
- Team BK: 12%
- Team KC: 4%
- Team IZI: 0%
- Unsupported tie break scenario: 477
