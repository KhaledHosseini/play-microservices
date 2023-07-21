from kafka import KafkaConsumer
import json
import logging
from typing import List


class KafkaConsumerWorker:
    def __init__(self, topic: str, consumerGroupId: str, kafkaBrokers: List[str]):
        self.Topic = topic
        self.ConsumerGroupId = consumerGroupId
        self.KafkaBrokers = kafkaBrokers

    def run_kafka_consumer(self):
        consumer = KafkaConsumer(self.Topic,
        group_id=self.ConsumerGroupId, 
        bootstrap_servers=self.KafkaBrokers,
        enable_auto_commit=False,
        value_deserializer= self.loadJson)

        logging.info("consumer is listening....")
        try:
            for message in consumer:
                json = self.loadJson(message.value)
                self.messageRecieved(self.Topic, json)
                #committing message manually after reading from the topic
                consumer.commit()
        except Exception as e:
            logging.error(f"consumer listener stoped with error: {e}")
        finally:
            consumer.close()

    def messageRecieved(self,topic: str,message: any):
        raise Exception("messageRecieved not implemented")
        

    def loadJson(self,value):
        logging.info(f"decoding message: {value}")
        try:
            if self.is_encoded_message(value):
                js = json.loads(value.decode('utf-8'))
                return js
            else:
                return value
        except json.decoder.JSONDecodeError as e: 
            logging.error("invalid json:", e)
            return "invalid json"
        
    def is_encoded_message(self,message):
        if isinstance(message, bytes):
            return True
        elif isinstance(message, str):
            try:
                message.encode('utf-8')
            except UnicodeEncodeError:
                return False
            else:
                return True
        else:
            return False
    def run(self):
        self.run_kafka_consumer()
