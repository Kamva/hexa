package main

type errLayer struct {
    app App
}

func (a *errLayer) By() (string,error) {
        r0,r1 := a.app.By()
        
            if r1 != nil {
            // do something when error occured.
            }
        
        return r0,r1
    
    }

func (a *errLayer) Hi(username string,password string) (sentence string,err error) {
        r0,r1 := a.app.Hi(username,password)
        
            if r1 != nil {
            // do something when error occured.
            }
        
        return r0,r1
    
    }

func NewErrLayer(app App) App {
    return &errLayer{
    app: app,
    }
}
