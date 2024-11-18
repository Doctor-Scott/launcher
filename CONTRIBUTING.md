If you'd like to contribute, create a pull request. You can also open an issue for any bugs or feature requests.

The repo has a `.githook/` folder for the pre-commit hooks, these run the tests before a commit

>This is a good idea to do if you don't want to break things too badly
```sh
git config core.hooksPath .githooks
chmod +x .githooks/pre-commit      
```

### Workflows
Contribute a saved workflow by adding the created `json` file to the `workflows/community` directory
