# water
Keep your plants watering schedule on track

# Generation tab completion script
Cobra based command can generate tab completion script.

For zsh, put either
```
source <(./water completion zsh)
```
or
```
./water completion zsh > "${fpath[1]}/_water"
autoload -U compinit && compinit
```
 to ~/.zshrc

# Build
```
go build
```
Or
```
go build -o ~/go/bin/water
```
to run the binary everywhere