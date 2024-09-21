FROM bonavadeur/gython:root

MAINTAINER bonavadeur

WORKDIR /

COPY ./outlier /ko-app/outlier/
COPY ./requirements.txt /ko-app/outlier/
COPY ./hack/startup.sh /ko-app/

RUN pip install -r /ko-app/outlier/requirements.txt

RUN chmod +x /ko-app/startup.sh

ENTRYPOINT ["/ko-app/startup.sh"]
