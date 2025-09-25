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

## 🚀 CI/CD Workflow

This project uses GitHub Actions for Continuous Integration and Continuous Deployment. The workflow is defined in [`./github/workflows/go-ci.yml`](.github/workflows/go-ci.yml).

### How CI/CD Works

The CI/CD pipeline is configured to build, test, and deploy the `GoTerminal` application based on pushes to specific branches:

*   **`dev_backend` branch:**
    *   Triggers the `build-and-test` job.
    *   Builds the Go application and runs unit tests.
    *   Builds a Docker image and pushes it to Docker Hub with tags: `your-dockerhub-username/goterminal-app:<commit_sha>` and `your-dockerhub-username/goterminal-app:dev_backend`.
    *   **Note:** There is no automatic deployment to Kubernetes for `dev_backend`. This branch is primarily for development and testing the Docker image build process.

*   **`staging` branch:**
    *   Triggers the `build-and-test` job, followed by the `deploy-staging` job.
    *   The `deploy-staging` job sets up Kubernetes context for the staging environment and applies the `kubernetes/deployment.yaml`, `kubernetes/service.yaml`, and `kubernetes/hpa.yaml` configurations, deploying the latest image built from the `staging` branch.

*   **`production` branch:**
    *   Triggers the `build-and-test` job, followed by the `deploy-production` job.
    *   The `deploy-production` job sets up Kubernetes context for the production environment and applies the `kubernetes/deployment.yaml`, `kubernetes/service.yaml`, and `kubernetes/hpa.yaml` configurations, deploying the latest image built from the `production` branch.

### How to See CI/CD and Deployment Work

1.  **Push to a Branch:**
    *   **For `dev_backend`:** Make a change and push to `dev_backend`. Monitor the "Actions" tab on GitHub to see the `build-and-test` job complete and the Docker image pushed. You can then manually pull and run the Docker image or deploy it to a local Kubernetes cluster.
    *   **For `staging` or `production`:** Make a change and push to `staging` or `production`. Monitor the "Actions" tab on GitHub. After the `build-and-test` job, the respective `deploy-staging` or `deploy-production` job will run, deploying the application to the Kubernetes cluster.

2.  **Verify Deployment in Kubernetes:**
    After a successful deployment to `staging` or `production`, you can verify the deployment using `kubectl`:
    ```bash
    kubectl get deployments
    kubectl get pods
    kubectl get services
    ```
    Ensure `goterminal-deployment` and `goterminal-service` are listed, and your `goterminal` pod is in a `Running` state with `1/1` ready.

3.  **Access the Web Application (Once Implemented):**
    Once `GoTerminal` is converted into a web application, you can access it via the Kubernetes service. If your Kubernetes cluster is running locally (e.g., Minikube or Docker Desktop), you can use `minikube service goterminal-service` (for Minikube) or `kubectl port-forward service/goterminal-service 8080:80` (for Docker Desktop) to get the URL or access it locally. For remote clusters, you would typically configure an Ingress controller or a LoadBalancer service to expose the application externally.

## 🔒 Environment Variables and GitHub Secrets

This project utilizes several environment variables and GitHub Secrets for secure and flexible CI/CD operations.

### Required Variables

The following variables need to be configured:

*   **`DOCKER_USERNAME`**: Your Docker Hub username.
*   **`DOCKER_PASSWORD`**: Your Docker Hub password or a Personal Access Token.
*   **`KUBE_CONFIG_STAGING`**: Base64 encoded content of your Kubernetes kubeconfig file for the staging environment.
*   **`KUBE_CONFIG_PROD`**: Base64 encoded content of your Kubernetes kubeconfig file for the production environment.
*   **`DOCKER_REPO`**: The full path to your Docker repository (e.g., `your-dockerhub-username/goterminal-app`).

### Configuration Steps

1.  **Local Development (`config.env.example`):**
    For local testing or manual Docker operations, you can create a `.env` file based on `config.env.example` in the root directory of this project.
    ```bash
    cp config.env.example .env
    # Edit .env with your actual values
    ```
    **Important:** Never commit your `.env` file to version control.

2.  **GitHub Secrets (for CI/CD):**
    For the GitHub Actions CI/CD pipeline to function correctly, you must configure these variables as GitHub Secrets in your repository.

    *   Go to your GitHub repository.
    *   Navigate to `Settings` > `Secrets and variables` > `Actions`.
    *   Click `New repository secret` and add each of the required variables (`DOCKER_USERNAME`, `DOCKER_PASSWORD`, `KUBE_CONFIG_STAGING`, `KUBE_CONFIG_PROD`, `DOCKER_REPO`) with their respective values.

    **Note on Kubeconfig:** To get the base64 encoded value of your kubeconfig, you can use the following command:
    ```bash
    cat ~/.kube/config | base64
    # On macOS, you might need: cat ~/.kube/config | base64 | tr -d '\n'
    ```
    Ensure you use the correct kubeconfig file for your staging and production clusters.

## 🛠 Installation & Usage

### Prerequisites

*   Go installed (version ≥ 1.16)
*   Docker installed and running
*   `kubectl` installed and configured

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

