# controller-ownref-experiment

```
controller-ownref-experiment on  main [?] via  v1.21.4 on   (eu-central-1)
❯ kubectl get secrets -n controller-ownref-experiment-system
No resources found in controller-ownref-experiment-system namespace.

controller-ownref-experiment on  main [?] via  v1.21.4 on   (eu-central-1)
❯ kubectl apply -f example/simple.yaml
myres.example.k8s.experiments.efertone.me/myres-example created

controller-ownref-experiment on  main [?] via  v1.21.4 on   (eu-central-1)
❯ kubectl get secrets -n controller-ownref-experiment-system
NAME                         TYPE     DATA   AGE
myres-example-fancy-secret   Opaque   0      2s

controller-ownref-experiment on  main [?] via  v1.21.4 on   (eu-central-1)
❯ kubectl delete -f example/simple.yaml
myres.example.k8s.experiments.efertone.me "myres-example" deleted

controller-ownref-experiment on  main [?] via  v1.21.4 on   (eu-central-1)
❯ kubectl get secrets -n controller-ownref-experiment-system
No resources found in controller-ownref-experiment-system namespace.
```
