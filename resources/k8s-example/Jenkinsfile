node {
    checkout scm

    environment {
       DOCKER_HUB_ACCOUNT = 'icrosby'
       DOCKER_IMAGE_NAME = 'k8s-example-adidas'
    }

    def DOCKER_HUB_ACCOUNT = 'icrosby'
    def DOCKER_IMAGE_NAME = 'k8s-example-adidas'

    echo 'Building Go App'
    stage("build") {
        docker.image("icrosby/jenkins-agent:kube").inside('-u root') {
            sh 'go build' 
        }
    }
    echo 'Testing Go App'
    stage("test") {
        docker.image('icrosby/jenkins-agent:kube').inside('-u root') {
            sh 'go test' 
        }
    }

    echo 'Building Docker image'
    stage('BuildImage')
    def app = docker.build("${DOCKER_HUB_ACCOUNT}/${DOCKER_IMAGE_NAME}:${JOB_NAME}_${BUILD_NUMBER}", '.')

    echo 'Testing Docker image'
    stage("test image") {
        docker.image("${DOCKER_HUB_ACCOUNT}/${DOCKER_IMAGE_NAME}").inside {
            sh './test.sh'
        }
    }

    echo 'Pushing Docker Image'
    stage("Push")
    docker.withRegistry('https://index.docker.io/v1/', 'docker-hub') {
        app.push()
    }

    echo 'Pushing Docker Image Locally'
    stage("Push Local")
    docker.withRegistry('http://localhost:5000/') {
        app.push()
    }
    
    echo 'Tag as Production'
    stage("Tag")
    docker.withRegistry('https://index.docker.io/v1/', 'docker-hub') {
        app.push("production")
    }

    echo "Deploying image"
    stage("Deploy") 
    docker.image('smesch/kubectl').inside{
        withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
            sh "kubectl --kubeconfig=$KUBECONFIG set image deployment/k8s-example k8s-example=icrosby/k8s-example-adidas:${JOB_NAME}_${BUILD_NUMBER}"
        }
    }
}
