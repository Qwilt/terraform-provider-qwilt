
#Create an empty resource to import into.

#After the import is complete, manually set the required attributes 
#in the resource based on the imported state.

#We recommend changing the site_id attribute to a reference 
#to the qwilt_cdn_site resource to achieve implicit dependency.


resource "qwilt_cdn_site_configuration" "example" {
}

# You can import the qwilt_cdn_site_configuration resource by 
# specifying just the site_id. 

    #For example: terraform import qwilt_cdn_site_configuration.example <site_id>
    
        # The process determines which saved site version (represented by the 
        # revisionId) to import based on the following conditions: 
        # - If there is an active published version, it is imported.
        # - If not, the most recently published version is imported. 
        # - If the site has never been published, the most recently saved 
        #   configuration version is imported.

# Alternatively, you can specify a particular revision_id of the site 
# configuration by adding a colon (:) and the revision_id.

terraform import qwilt_cdn_site_configuration.example <site_id>:<revision_id>
