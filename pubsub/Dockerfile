FROM gcr.io/google.com/cloudsdktool/cloud-sdk:367.0.0-emulators

RUN apt-get update && \
    apt-get install -y git python3-pip netcat && \
    git clone https://github.com/googleapis/python-pubsub.git

ENV PROJECT_ID=change-stream-pubsub \
    TOPIC_ID=change-event-topic \
    SUBSCRIPTION_ID=change-event-subscription \
    PUSH_ENDPOINT=http://host.docker.internal:3000

WORKDIR /python-pubsub/samples/snippets
RUN pip3 install virtualenv && \
    virtualenv env && \
    . env/bin/activate && \
    pip3 install -r requirements.txt

COPY ./entrypoint.sh ./
EXPOSE 8085
ENTRYPOINT ["./entrypoint.sh"]
