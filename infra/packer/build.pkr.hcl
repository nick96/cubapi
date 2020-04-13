build {
  sources = [
    "source.digitalocean.base_image"
  ]


  // Copy the docker-compose file over
  provisioner "file" {
    source      = "../docker-compose.yml"
    destination = "/opt/docker-compose.yml"
  }


  // Copy the bagheera service file over
  provisioner "file" {
    source      = "bagheera.service"
    destination = "/etc/systemd/system/bagheera.service"
  }


  // Install docker
  provisioner "shell" {
    inline = [
      "sudo apt-get update",
      "sudo apt-get install -y apt-transport-https ca-certificates",
      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
      "sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\"",
      "sudo apt-get update",
      "sudo apt-get install -y docker-ce docker-ce-cli containerd.io",
    ]
  }

  // Install docker-compose
  provisioner "shell" {
    inline = [
      "sudo curl -L \"https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)\" -o /usr/local/bin/docker-compose",
      "sudo chmod +x /usr/local/bin/docker-compose"
    ]
  }

  // Enable and start the docker service
  provisioner "shell" {
    inline = [
      "systemctl enable docker",
      "systemctl start docker",
    ]
  }


  // Enable the bagheera service
  provisioner "shell" {
    inline = [
      "systemctl enable bagheera",
    ]
  }



  // Pull down all the docker images to improve startup time
  provisioner "shell" {
    inline = [
      "docker-compose -f /opt/docker-compose.yml pull",
    ]
  }
}
