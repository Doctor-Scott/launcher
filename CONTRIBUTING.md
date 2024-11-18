# Contributing
If you'd like to contribute, create a pull request. You can also open an issue for any bugs or feature requests.

## Development helpers
### Pre-commit hooks
The repo has a `.githook/` folder for the pre-commit hooks, these run the tests before a commit

Run these commands in the root of the project to set these up:
>This is a good idea to do if you don't want to break things too badly
```sh
git config core.hooksPath .githooks
chmod +x .githooks/pre-commit      
```

### Live reload
There is also a live reload feature, which will automatically rebuild the launcher when you make changes to the code.

First, in one terminal:
```sh
./live_reload/rebuild.sh
```
Then, in another terminal:
```sh
./live_reload/watch.sh
```
This will run the launcher reloading when `rebuild` compiles the code:

# Workflows
Contribute a saved workflow by adding the created `json` file to the `workflows/community` directory
