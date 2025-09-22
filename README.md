# 🧮 GoTerminal

**GoTerminal** is a simple yet powerful interactive command-line interface (CLI) built with Go. It features a variety of general and arithmetic commands, colorful output for readability, command history, file I/O, configuration management, and support for piping.

---

## 🚀 Features

### ✅ General Commands

| Command          | Description                      |
| ---------------- | -------------------------------- |
| `help`           | Display all available commands   |
| `exit`           | Exit the terminal gracefully     |
| `history`        | Show history of entered commands |
| `clear`          | Clear command history            |
| `echo <message>` | Echo back the provided message   |


### ➕ Arithmetic Operations

| Command                  | Description                           |
| ------------------------ | ------------------------------------- |
| `arithmetic`             | List all arithmetic commands          |
| `add <num1> <num2>`      | Add two numbers                       |
| `subtract <num1> <num2>` | Subtract second number from the first |
| `multiply <num1> <num2>` | Multiply two numbers                  |
| `division <num1> <num2>` | Divide first number by the second     |
| `modulus <num1> <num2>`  | Return modulus of the two numbers     |

### 📂 File I/O

| Command                       | Description                        |
| ----------------------------- | ---------------------------------- |
| `cat <filename>`              | Display contents of a file         |
| `write <filename> <content>`  | Write content to a file            |
| `append <filename> <content>` | Append content to an existing file |
| `save`                        | Save to a file with options        |
| `save history`                | Save command history to `command_history.txt` |
| `save arithmetic`             | Save arithmetic operations to `arithmetic_operations.txt` |

### 🔒 Permissions and File Access
*   **Integrated Permission Checks**: Commands interacting with the file system (`cat`, `write`, `append`, `save history`, `save arithmetic`) now include robust permission checks.
*   **Access Control**: The terminal verifies read and write permissions before executing file operations, preventing unauthorized access and modifications.
*   **Error Handling**: Clear error messages are provided when file access is denied or files are not found.


### 🧩 Pipe Support

* Chain multiple commands using:
  `pipe <cmd1> | <cmd2> | ...`
  Passes output of one command as input to the next.

### 🎨 Colorful Output

* Terminal output is color-coded using [`fatih/color`](https://github.com/fatih/color) for:

  * Prompts
  * Errors
  * Standard output

### ⌨️ Readline Integration

Powered by [`chzyer/readline`](https://github.com/chzyer/readline):

* Line editing
* Arrow key history navigation
* Command autocompletion

### ⚙️ Configuration & Themes

* `config`: View current config
* `set <property> <value>`: Update config
* `theme <theme-name>`: Change terminal theme
  Available themes: `default`, `dark`, `light`, `solarized`

---


## 🛠 Installation & Usage

### Prerequisites

* Go installed (version ≥ 1.16)
* Docker installed and running
* `kubectl` installed and configured

### Install Dependencies

```bash
go get github.com/chzyer/readline
go get github.com/fatih/color
```

### Run the Application (Go Native)

```bash
go run main.go
```

### Run Unit Tests

To run the unit tests for GoTerminal, navigate to the project's root directory and execute the following command:

```bash
go test ./...
```

This command will discover and run all tests in the current module.

### Run the Application (Docker)

1.  **Build the Docker Image:**
    ```bash
    docker build -t goterminal-app .
    ```

2.  **Run the Docker Container (Interactive):**
    ```bash
    docker run -it goterminal-app ./goTerminal
    ```
    This will start the `goTerminal` application inside a Docker container, and you can interact with it directly in your terminal. Type `exit` to stop the container.

    **Note:** If you want to run the container in the background and then attach to it later, you can use:
    ```bash
    docker run -d --name goterminal-bg goterminal-app
    docker exec -it goterminal-bg ./goTerminal
    ```

### Deploy to Kubernetes

1.  **Build the Docker Image:**
    Build the Docker image for the `goTerminal` application:
    ```bash
    docker build -t goterminal-app:latest .
    ```

2.  **Load Docker Image into Kubernetes (for local clusters like Minikube/Docker Desktop):**
    If you are using a local Kubernetes cluster (like Minikube or Docker Desktop's Kubernetes), you need to ensure the newly built image is available within the cluster's Docker daemon.
    -   **For Minikube:**
        ```bash
        minikube image load goterminal-app:latest
        ```
    -   **For Docker Desktop Kubernetes:** The image built locally is usually already available to the Docker Desktop Kubernetes cluster. No extra step is typically needed.

    For remote Kubernetes clusters, push your Docker image to a container registry (e.g., Docker Hub) and update the `image` field in `kubernetes/deployment.yaml` to point to the registry path.

3.  **Apply Kubernetes Configurations:**
    Apply the Kubernetes deployment and service configurations. If a previous deployment exists, this command will update it.
    ```bash
    kubectl apply -f kubernetes/deployment.yaml
    kubectl apply -f kubernetes/service.yaml
    ```

4.  **Verify Deployment:**
    Check the status of your deployment and pods:
    ```bash
    kubectl get deployments
    kubectl get pods
    kubectl get services
    ```
    Ensure `goterminal-deployment` and `goterminal-service` are listed, and your `goterminal` pod is in a `Running` state with `1/1` ready.

5.  **Interact with the Application in Kubernetes:**
    To interact with the CLI application running inside the Kubernetes pod, you need to execute commands within the pod. First, get the exact name of your running pod:
    ```bash
    kubectl get pods
    ```
    Look for a pod name starting with `goterminal-deployment-` (e.g., `goterminal-deployment-686794fd9f-abcde`). Then, execute the `goTerminal` application inside that pod:
    ```bash
    kubectl attach -it <your-goterminal-pod-name> -c goterminal
    # OR
    kubectl exec -it <your-goterminal-pod-name> -- ./goTerminal
    ```
    Replace `<your-goterminal-pod-name>` with the actual name you found. Once inside, you can use all `goTerminal` commands. Type `exit` to leave the pod's terminal.

---



## 🤝 Contributing

If you have suggestions for improvements, new features, or bug fixes, please follow these steps:

1.  **Fork the Repository:** Start by forking the GoTerminal repository to your GitHub account.
2.  **Clone Your Fork:** Clone your forked repository to your local machine:
    ```bash
    git clone https://github.com/MaxV16/GoTerminal.git
    cd GoTerminal
    ```
3.  **Create a New Branch:** Create a new branch for your feature or bug fix:
    ```bash
    git checkout -b feature/your-feature-name
    # or
    git checkout -b bugfix/issue-description
    ```
4.  **Make Your Changes:** Implement your changes, ensuring they adhere to the existing code style and conventions.
5.  **Test Your Changes:** If applicable, add or update tests to cover your modifications. Ensure all existing tests pass.
6.  **Commit Your Changes:** Commit your changes with a clear and descriptive commit message:
    ```bash
    git commit -m "feat: Add new feature"
    # or
    git commit -m "fix: Resolve issue with command history"
    ```
7.  **Push to Your Fork:** Push your new branch to your forked repository on GitHub:
    ```bash
    git push origin feature/your-feature-name
    ```
8.  **Open a Pull Request:** Go to the original GoTerminal repository on GitHub and open a new pull request from your forked branch. Provide a detailed description of your changes and why they are necessary.

