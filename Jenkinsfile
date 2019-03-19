pipeline {
    agent any
    stages {
        stage('Prepare') {
            steps {
                echo 'Prepare..'
            }
        }
        stage('Test') {
            steps {
                echo "Test.."
            }
        }
        stage("Analysis") {
            steps {
                echo 'Analysis..'
            }
        }
        stage("Check") {
            steps {
                echo 'Check..'
            }
        }
        stage('Dispatch') {
        	parallel {
				stage('OPS') {
					stages {
						stage('Adapter') {
							steps {
								echo "Adapter.."
							}
						}
						stage('Build') {
							steps {
								echo "Build.."
							}
						}
						stage('Docker') {
							steps {
								echo "Docker.."
							}
						}
						stage('Deploy') {
							steps {
								echo "Deploy.."
							}
						}
					}
                }
                stage('ZSY') {
					stages {
						stage('Adapter') {
							steps {
								echo "Adapter.."
							}
						}
						stage('Build') {
							steps {
								echo "Build.."
							}
						}
						stage('Docker') {
							steps {
								echo "Docker.."
							}
						}
						stage('Deploy') {
							steps {
								echo "Deploy.."
							}
						}
					}
                }
			}
        }
        stage("Recovery") {
            steps {
                echo 'Recovery..'
            }
        }
    }
}