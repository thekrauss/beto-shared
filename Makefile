DOCS_OUT=docs

.PHONY: docs
docs:
	mkdir -p $(DOCS_OUT)
	gomarkdoc ./pkg/errors > $(DOCS_OUT)/errors.md
	gomarkdoc ./pkg/gormhelpers > $(DOCS_OUT)/gormhelpers.md
	gomarkdoc ./pkg/storage > $(DOCS_OUT)/storage.md
	gomarkdoc ./pkg/openstack-client > $(DOCS_OUT)/openstack-client.md
	gomarkdoc ./pkg/redis > $(DOCS_OUT)/redis.md
	gomarkdoc ./pkg/metrics > $(DOCS_OUT)/metrics.md
	gomarkdoc ./pkg/middleware > $(DOCS_OUT)/middleware.md
	gomarkdoc ./pkg/eventbus > $(DOCS_OUT)/eventbus.md
	gomarkdoc ./pkg/authz > $(DOCS_OUT)/authz.md
	gomarkdoc ./pkg/tracing > $(DOCS_OUT)/tracing.md
