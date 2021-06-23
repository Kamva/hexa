### Distributed locks


#### Available drivers:
- [x] MongoDB
- [ ] Redis red locks.
- [ ] Etcd.

__Note__: currently Mongodb is good for us, but if someday we needed 
to the redis red locks or etcd, we can implement them.


#### Docs
You can lock everything between multiple instances of your app
using distributed locks. for example if you want to run a DB migration
using at you app's bootstrap time, to prevent for multiple run of that 
migration by each instance of the app on your server(say pods in a k8s cluster, e.g., `pod-my-app-893yr`)
you need to lock migration on one instance and other instances need to wait
for that lock to release, after that other instances can check and if migration 
has been done by another instance, so they skip it and run the app. 

We have two important concepts:
- __DLM__(distributed lock manager): manges locks in our apps.
- __Mutex__: Lock instance. We create mutex by `DLM`.
Each mutex has three fields:
    - `name` (`key`): the lock name.
    - `owner`: owner (e.g., machine that acquires the lock) of the lock.
    - `ttl`: The lock ttl.
    


#### Rules:
- Different mutexes with same `lock name` and `same
   owner` can lock and unlock each other. 
  
#### How to use?
Just create a new instance of DLM in your app as a service and everywhere
that you needed you can create a new `mutex` using your `dlm`.
Please note at `dlm` creation time, use your host name as owner of the lock
If you want to have different owners for each instance of the app.  
To see Examples check test files please.
  

