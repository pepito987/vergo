# vergo

# Simple usages
  * returns the latest tag/release prefixed with banana

    `vergo get --tag-prefix=banana`

  * increments patch part of the version prefixed with banana
      
    `vergo bump patch --tag-prefix=banana`
      
  * increments minor part of the version prefixed with banana
    
    `vergo bump minor --tag-prefix=banana`
  
  * increments minor part of the version prefixed with banana and pushes it to origin remote
  
    `vergo bump major --tag-prefix=apple --push-tag`