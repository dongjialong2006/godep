pipeline('godep1') {
    agent any

    stages {
        stage("Prepare") {
            steps {
                echo 'Prepare..'
            }
        }
        stage('Build') {
            steps {
                sh 'make'
                archiveArtifacts artifacts: 'bin/*', fingerprint: true
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}

pipeline('godep2') {
    agent any

    stages {
        stage('Build1') {
            steps {
                sh 'make'
                archiveArtifacts artifacts: 'bin/*', fingerprint: true
            }
        }
        stage('Test1') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy1') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}