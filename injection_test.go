package main

import (
	// "fmt"
	"testing"
)

const jsonText = `
{
  "ignition": {
    "version": "3.3.0"
  },
  "storage": {
    "directories": [
      {
        "path": "/tmp/example",
        "mode": 493
      },
      {
        "path": "/tmp/example2"
      },
      {
        "path": "/tmp/example3/example4",
        "mode": 484
      }
    ],
    "files": [
      {
        "path": "/tmp/example/hello_world.txt",
        "contents": {
          "compression": "",
          "source": "data:,Hello%2C%20world!%0A"
        }
      },
      {
        "path": "/tmp/example/hello_world_gzip.txt",
        "contents": {
          "compression": "gzip",
          "source": "data:;base64,H4sIAAAAAAAC/zSQQc5TMQyE9+8Uc4CqpwAJJJYg1iZxX0dK4tSxS8XpUf7CLlLs8Xzf9zsXuCDoWe5o4qc6vmhrhp/mreLGpoi7BH6zNeRSnH84UaxP17Vo4+BA3BVfz8GgjSu+mWsH58qOas0ciwHpGhcUG0tLaKRDKicXC8d5aGNcsLSiGpS5ulWE9mkOjsLKmiOQgSa/zBUa72xFl3PIIY2PlCt+BHSwQyo69+Opg9IveCQXhq3wrNCXemHILoxsTXqx4yN5D3FxX/qI5IS+oLKRu1V7EzxS4opPO1IyFPR0fRfaPlyn611HVWeAA09rOUNC8dyk0LUUha39V6TQxC1PSmDsQscUp0T6FZ9fRWdobo8jYKWIFgmUnKwSe8MGphurjm1xm+I4SrYpmxt2u7FQUHWp799ubdeQLYgVuv55zX49/gYAAP//1qfI0RYCAAA="
        },
        "mode": 420
      }
    ]
  },
  "systemd": {
    "units": [
      {
        "contents": "[Unit]\nDescription=Hello world service\n\n[Service]\nType=oneshot\nExecStart=/usr/bin/echo \"hello world\"\nStandardOutput=journal\n\n[Install]\nWantedBy=multi-user.target default.target\n",
        "enabled": true,
        "name": "hello-world.service"
      }
    ]
  }
}`
