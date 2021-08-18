pipeline {
    agent any

    options {
        disableConcurrentBuilds() //禁止同时执行
        buildDiscarder(logRotator(daysToKeepStr: '7', numToKeepStr: '10')) // 构建记录保存7天,最多保存10个构建记录
        skipDefaultCheckout()
        timeout(time: 1, unit: 'HOURS')
        timestamps()
    }

    parameters {
        choice(name: "mode", choices: ["编译发布","重启","回滚"], description: '<h5 style="color:#DA70D6;">选择模式</h5>')
    }


    environment {

        //环境配置
        user_tools_dir = "${JENKINS_HOME}/user_tools"
        PATH = "${user_tools_dir}/bin:$PATH"

        //获取基本参数
        app_name = "${JOB_BASE_NAME}"
        project_group = sh(script: "echo ${JOB_NAME}|awk -F'/' '{printf \$(NF-1)}'", returnStdout: true).trim()
        deploy_env = sh(script: "echo ${JOB_NAME}|awk -F'/' '{printf \$(NF-2)}'", returnStdout: true).trim()

        //钉钉消息通知开关
        dingtalk_enable = "true"
        
    }
    stages {
        stage ("获取代码") {
            when {
                environment name: 'mode', value: '编译发布'
            }
            steps {
                script {
                    GitBranch = pub.GetGit()
                    CommitId  = pub.GetCommitId()
                }
            }
        }

        stage ("配置环境") {
            steps {
                script {
                    (RegistryUrl,RegistryNS) = pub.RegistryLogin()
                    K8sNS = pub.GetKubeConf()
                    Version = pub.GetVersion()

                    DeployVersion = pub.CheckGray()
                    pub.SetJobName()

                    AppType = pub.GetAppType()
                    AppMem = pub.GetAppMem()
                }
            }
        }

        stage ("编译") {
            when {
                environment name: 'mode', value: '编译发布'
            }
            steps {
                script {
                    if (AppType == "java") {
                        (JarName,JarFile) = pub.BuildMaven()  
                    } else if (AppType == "node") {
                        pub.BuildNode() 
                    } else {
                        echo "AppType Error"
                        sh "exit 1"
                    }
                }
            }
        }

        stage ("镜像") {
            when {
                environment name: 'mode', value: '编译发布'
            }
            steps {
                script {
                    if (env.gray_deploy == "false") {
                        pub.BuildProImage()
                    } else if (AppType == "java") {
                        pub.BuildMavenImage() 
                    } else if (AppType == "node") {
                        pub.BuildNodeImage()
                    } else {
                        echo "AppType Error"
                        sh "exit 1"
                    }
                }
            }
        }

        stage ("部署") {
            when {
                environment name: 'mode', value: '编译发布'
            }
            steps {
                script {
                    pub.K8sDeploy()
                }
            }
        }

        stage ("重启") {
            when {
                environment name: 'mode', value: '重启'
            }
            steps {
                script {
                    pub.K8sRestart()
                }
            }
        }

        stage ("回滚") {
            when {
                environment name: 'mode', value: '回滚'
            }
            steps {
                script {
                    pub.K8sRollback()
                }
            }
        }

    }
    post {
        always {
            script {
                pub.DingTalk()
            }
        }
        cleanup {
            deleteDir()
            dir("${workspace}@tmp") {deleteDir()}
            dir("${workspace}@script") {deleteDir()}
        }
    }
}
